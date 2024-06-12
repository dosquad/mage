package mage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
	"github.com/na4ma4/go-permbits"
)

// Mirrord namespace is defined to group Mirrord functions.
type Mirrord mg.Namespace

// VsCodeDebugConfig generates debug launch.json in service project dir.
func (Mirrord) VsCodeDebugConfig(_ context.Context) error {
	launchCfg := `{
		// Use IntelliSense to learn about possible attributes.
		// Hover to view descriptions of existing attributes.
		// For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
		"version": "0.2.0",
		"configurations": [
		%s
		],
	}
	`

	launchItemCfg := `{
				"name": "Launch Service: %s",
				"type": "go",
				"request": "launch",
				"mode": "debug",
				"program": "${workspaceFolder}/cmd/%s",
				"hideSystemGoroutines": true,
				"envFile": "${workspaceFolder}/artifacts/data/.env",
				"env": {
					"MIRRORD_CONFIG_FILE": "${workspaceFolder}/mirrord.yaml"
				}
			}
`

	launchItems := []string{}
	for _, cmdPath := range helper.MustCommandPaths() {
		cmdPath = filepath.Base(cmdPath)
		launchItems = append(launchItems,
			fmt.Sprintf(launchItemCfg, cmdPath, cmdPath),
		)
	}

	launchBody := fmt.Sprintf(launchCfg, strings.Join(launchItems, ","))

	helper.MustMakeDir(helper.MustGetVSCodePath(), permbits.MustString("ug=rwx,o=rx"))

	cfgFilePath := filepath.Join(helper.MustGetVSCodePath(), "launch.json")

	f, err := os.Create(cfgFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(launchBody)
	return err
}
