package mage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/dosquad/mage/helper"
	"github.com/magefile/mage/mg"
)

// Mirrord namespace is defined to group Mirrord functions.
type Mirrord mg.Namespace

// VsCodeDebugConfig generates debug launch.json in service project dir.
func (Mirrord) VsCodeDebugConfig(_ context.Context) error {
	launchCfg := fmt.Sprintf(`{
		// Use IntelliSense to learn about possible attributes.
		// Hover to view descriptions of existing attributes.
		// For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
		"version": "0.2.0",
		"configurations": [
			{	
				"name": "Launch Service",
				"type": "go",
				"request": "launch",
				"mode": "exec",
				"program": "${workspaceFolder}/artifacts/build/debug/darwin/arm64/${workspaceFolderBasename}",
				"hideSystemGoroutines": true,
				"env":{
					"MIRRORD_CONFIG_FILE": "%s/artifacts/config/mirrord.yaml"
				}
			}
		],
	}
	`, helper.MustGetWD())

	cfgFilePath := filepath.Join(helper.MustGetVSCodePath(), "launch.json")

	f, err := os.Create(cfgFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(launchCfg)
	return err
}
