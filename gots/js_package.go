package gots

import (
	"encoding/json"
	"path"
)

type PackageConfig struct {
	Main string `json:"main"`
}

func extractPackageConfig(tsModuleDir string) (*PackageConfig, string, error) {
	pkgJsonFilePath := path.Join(tsModuleDir, "package.json")
	bs, err := ReadFile(pkgJsonFilePath)
	if err != nil {
		return nil, "", err
	}
	pkgConfig := &PackageConfig{}
	err = json.Unmarshal(bs, pkgConfig)
	if err != nil {
		return nil, "", err
	}
	absDirPath, err := GetAbsPath(tsModuleDir)
	return pkgConfig, absDirPath, nil
}
