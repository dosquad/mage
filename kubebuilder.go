package mage

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
	"github.com/princjef/mageutil/shellcmd"
)

// Kubebuilder namespace is defined to group Kubebuilder functions.
type Kubebuilder mg.Namespace

// Manifests Generate WebhookConfiguration, ClusterRole and CustomResourceDefinition objects.
func (Kubebuilder) Manifests(_ context.Context) error {
	_ = helper.BinKubeControllerGen().Ensure()

	return helper.BinKubeControllerGen().Command(
		`rbac:roleName=manager-role crd webhook paths="./..." output:crd:artifacts:config=config/crd/bases`,
	).Run()
}

// Generate Generate code containing DeepCopy, DeepCopyInto, and DeepCopyObject method implementations.
func (Kubebuilder) Generate(_ context.Context) error {
	_ = helper.BinKubeControllerGen().Ensure()

	return helper.BinKubeControllerGen().Command(`object:headerFile="hack/boilerplate.go.txt" paths="./..."`).Run()
}

func kustomizeBuildCommand(cmd, filename string) (string, error) {
	out, err := helper.Command(string(helper.BinKustomize().Command(cmd)))
	if err != nil {
		return "", err
	}

	k8sPath := helper.MustGetArtifactPath("k8s")
	helper.MustMakeDir(
		k8sPath,
		permbits.MustString("u=rwx,go=rx"),
	)

	return filepath.Join(k8sPath, filename), helper.FileWrite(
		out,
		filepath.Join(k8sPath, filename),
	)
}

// Install CRDs into the K8s cluster specified in ~/.kube/config.
func (Kubebuilder) Install(ctx context.Context) error {
	mg.CtxDeps(ctx, Kubebuilder.Manifests)

	var crdManifest string
	{
		var err error
		if crdManifest, err = kustomizeBuildCommand("build config/crd", "crd.yaml"); err != nil {
			return err
		}
	}

	if err := shellcmd.Command(fmt.Sprintf(`kubectl apply -f "%s"`, crdManifest)).Run(); err != nil {
		return err
	}

	return nil
}

// Uninstall CRDs from the K8s cluster specified in ~/.kube/config.
// Call with ignore-not-found=true to ignore resource not found errors during deletion.
func (Kubebuilder) Uninstall(ctx context.Context) error {
	mg.CtxDeps(ctx, Kubebuilder.Manifests)

	var crdManifest string
	{
		var err error
		if crdManifest, err = kustomizeBuildCommand("build config/crd", "crd.yaml"); err != nil {
			return err
		}
	}

	if err := shellcmd.Command(
		fmt.Sprintf(`kubectl delete --ignore-not-found=true -f "%s"`, crdManifest),
	).Run(); err != nil {
		return err
	}

	return nil
}

// .PHONY: deploy
// deploy: manifests $(KUSTOMIZE) ## Deploy controller to the K8s cluster specified in ~/.kube/config.
// 	cd config/manager && ../../$(KUSTOMIZE) edit set image controller=${IMG}
// 	$(KUSTOMIZE) build config/default | kubectl apply -f -

// Deploy controller to the K8s cluster specified in ~/.kube/config.
func (Kubebuilder) Deploy(ctx context.Context) error {
	mg.CtxDeps(ctx, Kubebuilder.Manifests)

	_ = helper.BinKustomize()

	var dcfg *helper.DockerConfig
	{
		var err error
		dcfg, err = helper.DockerLoadConfig()
		helper.PanicIfError(err, "unable to load docker config")
	}

	{
		err := shellcmd.Command(
			fmt.Sprintf(`cd %s; %s edit set image controller=%s`,
				helper.MustGetGitTopLevel("config", "manager"),
				helper.MustGetArtifactPath("bin", "kustomize"),
				dcfg.GetImageRef(),
			),
		).Run()
		if err != nil {
			return err
		}
	}

	var deployManifest string
	{
		var err error
		if deployManifest, err = kustomizeBuildCommand("build config/default", "deploy.yaml"); err != nil {
			return err
		}
	}

	if err := shellcmd.Command(fmt.Sprintf(`kubectl apply -f "%s"`, deployManifest)).Run(); err != nil {
		return err
	}

	return nil
}

// Undeploy Undeploy controller from the K8s cluster specified in ~/.kube/config.
// Call with ignore-not-found=true to ignore resource not found errors during deletion.
func (Kubebuilder) Undeploy(ctx context.Context) error {
	mg.CtxDeps(ctx, Kubebuilder.Manifests)

	var deployManifest string
	{
		var err error
		if deployManifest, err = kustomizeBuildCommand("build config/default", "deploy.yaml"); err != nil {
			return err
		}
	}

	if err := shellcmd.Command(
		fmt.Sprintf(`kubectl delete --ignore-not-found=true -f "%s"`, deployManifest),
	).Run(); err != nil {
		return err
	}

	return nil
}

// Run a controller from your host.
func (Kubebuilder) Run(ctx context.Context) error {
	mg.CtxDeps(ctx, Kubebuilder.Manifests)
	mg.CtxDeps(ctx, Kubebuilder.Generate)
	mg.CtxDeps(ctx, Golang.Fmt)
	mg.CtxDeps(ctx, Golang.Vet)

	return shellcmd.Command("go run ./cmd/main.go --config=artifacts/data/config.yml -zap-devel").Run()
}
