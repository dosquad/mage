package mage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

// CFSSL namespace is defined to group CFSSL functions.
type CFSSL mg.Namespace

func mustCertDir(path ...string) string {
	return helper.MustGetArtifactPath(append([]string{"certs"}, path...)...)
}

// Install CFSSL binaries.
func (CFSSL) Install(_ context.Context) error {
	if err := helper.BinCfssl().Ensure(); err != nil {
		panic(err)
	}

	if err := helper.BinCfsslJSON().Ensure(); err != nil {
		panic(err)
	}

	return nil
}

func cfsslGenCert(_ context.Context, outputFile, configFileName, profile, srcFileName string) error {
	helper.PrintDebug("cfsslGenCert(ctx, %s, %s, %s, %s)", outputFile, configFileName, profile, srcFileName)
	var initCA []byte
	{
		var err error
		initCA, err = helper.Command(string(helper.BinCfssl().Command(fmt.Sprintf(
			`gencert -initca -config="%s" -profile="%s" "%s"`,
			configFileName, profile, srcFileName,
		))))
		if err != nil {
			return fmt.Errorf("unable to generate certificate: %w", err)
		}
	}

	var initCAf *os.File
	{
		var err error
		initCAf, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("unable to write certificate generate output to file(%s): %w", outputFile, err)
		}
		defer initCAf.Close()
	}

	if _, err := io.Copy(initCAf, bytes.NewReader(initCA)); err != nil {
		return fmt.Errorf("unable to stream certificate generate output to file(%s): %w", outputFile, err)
	}

	return nil
}

func cfsslSignCert(_ context.Context, outputFile, configFileName, profile, srcFileName, baseName string) error {
	helper.PrintDebug("cfsslSignCert(ctx, %s, %s, %s, %s)", outputFile, configFileName, profile, srcFileName)
	var initCA []byte
	{
		var err error
		initCA, err = helper.Command(string(helper.BinCfssl().Command(fmt.Sprintf(
			`sign -ca="%s" -ca-key="%s" -config="%s" -profile="%s" -csr="%s" "%s"`,
			mustCertDir("ca.pem"),
			mustCertDir("ca-key.pem"),
			configFileName, profile,
			baseName+".csr",
			srcFileName,
		))))
		if err != nil {
			return fmt.Errorf("unable to sign certificate: %w", err)
		}
	}

	var initCAf *os.File
	{
		var err error
		initCAf, err = os.Create(outputFile)
		if err != nil {
			return fmt.Errorf("unable to write certificate signature output to file(%s): %w", outputFile, err)
		}
		defer initCAf.Close()
	}

	if _, err := io.Copy(initCAf, bytes.NewReader(initCA)); err != nil {
		return fmt.Errorf("unable to stream certificate signature output to file(%s): %w", outputFile, err)
	}

	return nil
}

func cfsslJSON(_ context.Context, outputBase, inputFile string) error {
	helper.PrintDebug("cfsslJSON(ctx, %s, %s)", outputBase, inputFile)

	err := helper.BinCfsslJSON().Command(fmt.Sprintf(
		`-f %s -bare %s`,
		inputFile, outputBase,
	)).Run()
	if err != nil {
		return err
	}

	return nil
}

func cfsslInitCA(ctx context.Context) error {
	helper.PrintDebug("cfsslInitCA(ctx)")
	helper.PrintInfo("Generating and signing CA certificate")

	if helper.FileExists(mustCertDir("ca.pem")) {
		helper.PrintDebug("Target exists, skipping: %s", mustCertDir("ca.pem"))
		return nil
	}
	helper.MustMakeDir(mustCertDir(), permbits.MustString("ug=rwx,o=rx"))

	// Init CA
	if err := cfsslGenCert(
		ctx,
		mustCertDir("ca.json"),
		mustCertDir("ca-config.json"),
		"ca",
		helper.MustGetWD("testdata", "ca-csr.json"),
	); err != nil {
		return err
	}
	if err := cfsslJSON(
		ctx,
		mustCertDir("ca"),
		mustCertDir("ca.json"),
	); err != nil {
		return err
	}

	// Sign CA
	if err := cfsslSignCert(
		ctx,
		mustCertDir("ca-sign.json"),
		mustCertDir("ca-config.json"),
		"ca",
		helper.MustGetWD("testdata", "ca-csr.json"),
		mustCertDir("ca"),
	); err != nil {
		return err
	}
	if err := cfsslJSON(
		ctx,
		mustCertDir("ca"),
		mustCertDir("ca-sign.json"),
	); err != nil {
		return err
	}

	return nil
}

func cfsslCert(ctx context.Context, profile, baseName string) error {
	helper.PrintDebug("cfsslCert(ctx, %s, %s)", profile, baseName)
	helper.PrintInfo("Generating and signing certificate: %s", profile)

	if helper.FileExists(mustCertDir(baseName + ".pem")) {
		helper.PrintDebug("Target exists, skipping: %s", mustCertDir(baseName+".pem"))
		return nil
	}
	helper.MustMakeDir(mustCertDir(), permbits.MustString("ug=rwx,o=rx"))

	// Generate Cert
	if err := cfsslGenCert(
		ctx,
		mustCertDir(baseName+".json"),
		mustCertDir("ca-config.json"),
		profile,
		helper.MustGetWD("testdata", baseName+".json"),
	); err != nil {
		return err
	}
	if err := cfsslJSON(
		ctx,
		mustCertDir(baseName),
		mustCertDir(baseName+".json"),
	); err != nil {
		return err
	}

	// Sign Cert
	if err := cfsslSignCert(
		ctx,
		mustCertDir(baseName+"-sign.json"),
		mustCertDir("ca-config.json"),
		profile,
		helper.MustGetWD("testdata", baseName+".json"),
		mustCertDir(baseName),
	); err != nil {
		return err
	}
	if err := cfsslJSON(
		ctx,
		mustCertDir(baseName),
		mustCertDir(baseName+"-sign.json"),
	); err != nil {
		return err
	}

	return nil
}

func (CFSSL) Generate(ctx context.Context) error {
	if err := helper.BinInstall(helper.BinCfssl()); err != nil {
		return err
	}
	if err := helper.BinInstall(helper.BinCfsslJSON()); err != nil {
		return err
	}

	caConfig := helper.MustGetWD("testdata", "ca-config.json")
	if !helper.FileExists(caConfig) {
		return errors.New("testdata/ca-config.json not found")
	}

	helper.MustMakeDir(
		helper.MustGetArtifactPath("certs"),
		permbits.MustString("a=rx,ug=w"),
	)

	if err := helper.FileCopy(
		"testdata/ca-config.json",
		"artifacts/certs/ca-config.json",
		false,
	); err != nil && !errors.Is(err, helper.ErrFileExists) {
		return fmt.Errorf("unable to copy testdata/ca-config.json to artifacts/certs/ca-config.json: %w", err)
	}

	if err := cfsslInitCA(ctx); err != nil {
		return err
	}

	if helper.FileExists(helper.MustGetWD("testdata", "server.json")) {
		if err := cfsslCert(ctx, "server", "server"); err != nil {
			return err
		}
	}

	if helper.FileExists(helper.MustGetWD("testdata", "client.json")) {
		if err := cfsslCert(ctx, "client", "client"); err != nil {
			return err
		}
	}

	if helper.FileExists(helper.MustGetWD("testdata", "cert.json")) {
		if err := cfsslCert(ctx, "client", "cert"); err != nil {
			return err
		}
		if !helper.FileExists(mustCertDir("key.pem")) {
			if err := helper.FileCopy(
				mustCertDir("cert-key.pem"),
				mustCertDir("key.pem"),
				false,
			); err != nil {
				return err
			}
		}
	}

	return nil
}
