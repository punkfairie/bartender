package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type YamlStructure struct {
	InstallerPreferences InstallerPreferences `yaml:"installerPreferences"`
	SoftwarePackages     SoftwarePackages     `yaml:"softwarePackages"`
}

type InstallerPreferences struct {
	Apt     []string `yaml:"apt"`
	Darwin  []string `yaml:"darwin"`
	Dnf     []string `yaml:"dnf"`
	Freebsd []string `yaml:"freebsd"`
	Pacman  []string `yaml:"pacman"`
	Ubuntu  []string `yaml:"ubuntu"`
	Windows []string `yaml:"windows"`
	Zypper  []string `yaml:"zypper"`
}

type SoftwarePackages map[string]SoftwareDef

type SoftwareDef struct {
	Bin        *string   `yaml:"_bin"`
	Desc       string    `yaml:"_desc"`
	Deps       *[]string `yaml:"_deps"`
	Docs       *string   `yaml:"_docs"`
	Github     *string   `yaml:"_github"`
	Home       *string   `yaml:"_home"`
	Name       string    `yaml:"_name"`
	Apk        *string   `yaml:"apk"`
	Appimage   *string   `yaml:"appimage"`
	Basher     *string   `yaml:"basher"`
	Binary     *osNames  `yaml:"binary"`
	Bpkg       *string   `yaml:"bpkg"`
	Brew       *string   `yaml:"brew"`
	Cargo      *string   `yaml:"cargo"`
	Cask       *string   `yaml:"cask"`
	Crew       *string   `yaml:"crew"`
	Choco      *string   `yaml:"choco"`
	Dnf        *string   `yaml:"dnf"`
	Flatpak    *string   `yaml:"flatpak"`
	Gem        *string   `yaml:"gem"`
	Go         *string   `yaml:"go"`
	Krew       *string   `yaml:"krew"`
	Nix        *string   `yaml:"nix"`
	Npm        *string   `yaml:"npm"`
	Pacman     *string   `yaml:"pacman"`
	Pipx       *string   `yaml:"pipx"`
	PkgFreebsd *string   `yaml:"pkg-freebsd"`
	PkgTermux  *string   `yaml:"pkg-termux"`
	Port       *string   `yaml:"port"`
	Scoop      *string   `yaml:"scoop"`
	Script     *string   `yaml:"string"`
	Snap       *string   `yaml:"snap"`
	Whalebrew  *string   `yaml:"whalebrew"`
	Winget     *string   `yaml:"winget"`
	Xbps       *string   `yaml:"xbps"`
	Yay        *string   `yaml:"yay"`
	Zypper     *string   `yaml:"zypper"`
}

type osNames struct {
	Darwin  *string `yaml:"darwin"`
	Linux   *string `yaml:"linux"`
	Windows *string `yaml:"windows"`
}

func readYaml(file string) tea.Cmd {
	return func() tea.Msg {
		fileData, fileErr := os.ReadFile(file)
		if fileErr != nil {
			return errMsg{fileErr}
		}

		var parsedYaml YamlStructure

		yamlErr := yaml.Unmarshal(fileData, &parsedYaml)
		if yamlErr != nil {
			return errMsg{yamlErr}
		}

		return yamlMsg(parsedYaml)
	}
}
