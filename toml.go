package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

type YamlStructure struct {
	InstallerPreference InstallerPreference `toml:"installerPreference"`
	SoftwarePackages    SoftwarePackages    `toml:"softwarePackages"`
}

type InstallerPreference struct {
	Apt      []string `toml:"apt"`
	Darwin   []string `toml:"darwin"`
	Fedora   []string `toml:"fedora"`
	Freebsd  []string `toml:"freebsd"`
	Arch     []string `toml:"arch"`
	Ubuntu   []string `toml:"ubuntu"`
	Windows  []string `toml:"windows"`
	OpenSUSE []string `toml:"openSUSE"`
}

type SoftwarePackages map[string]SoftwareDef

type SoftwareDef struct {
	App        string            `toml:"_app"`
	Bin        map[string]string `toml:"_bin"`
	Deprecated bool              `toml:"_deprecated"`
	Desc       string            `toml:"_desc"`
	Deps       []string          `toml:"_deps"`
	Name       string            `toml:"_name"`
	Post       map[string]string `toml:"_post"`
	Service    string            `toml:"_service"`
	Systemd    string            `toml:"_systemd"`
	When       map[string]string `toml:"_when"`
	Apk        pkg               `toml:"apk"`
	Appimage   string            `toml:"appimage"`
	Apt        pkg               `toml:"apt"`
	Basher     pkg               `toml:"basher"`
	Binary     osNames           `toml:"binary"`
	Bpkg       pkg               `toml:"bpkg"`
	Brew       pkg               `toml:"brew"`
	BrewDarwin pkg               `toml:"brew:darwin"`
	Cargo      pkg               `toml:"cargo"`
	Cask       pkg               `toml:"cask"`
	Crew       pkg               `toml:"crew"`
	Choco      pkg               `toml:"choco"`
	Dnf        pkg               `toml:"dnf"`
	Flatpak    pkg               `toml:"flatpak"`
	Gem        pkg               `toml:"gem"`
	Go         pkg               `toml:"go"`
	Krew       pkg               `toml:"krew"`
	Nix        pkg               `toml:"nix"`
	Npm        pkg               `toml:"npm"`
	Pacman     pkg               `toml:"pacman"`
	Pipx       pkg               `toml:"pipx"`
	PkgFreebsd pkg               `toml:"pkg-freebsd"`
	PkgTermux  pkg               `toml:"pkg-termux"`
	Port       pkg               `toml:"port"`
	Scoop      pkg               `toml:"scoop"`
	Script     osNames           `toml:"script"`
	Snap       pkg               `toml:"snap"`
	Whalebrew  pkg               `toml:"whalebrew"`
	Winget     pkg               `toml:"winget"`
	Xbps       pkg               `toml:"xbps"`
	Yay        pkg               `toml:"yay"`
	Zypper     pkg               `toml:"zypper"`
}

type pkg interface {
	String() string
	Slice() []string
	Length() int
}

type singlePkg string

func (s singlePkg) String() string {
	return string(s)
}

func (s singlePkg) Slice() []string {
	return []string{string(s)}
}

func (s singlePkg) Length() int {
	return 1
}

type multiPkg []string

func (m multiPkg) String() string {
	return m.String()
}

func (m multiPkg) Slice() []string {
	return m
}

func (m multiPkg) Length() int {
	return len(m)
}

type osNames struct {
	Darwin  *string `toml:"darwin"`
	Linux   *string `toml:"linux"`
	Windows *string `toml:"windows"`
}

func (sd SoftwarePackages) UnmarshalTOML(data any) error {
	d, _ := data.(map[string]any)
	defs, _ := d["softwarePackages"].(map[string]map[string]any)

	for pkg, def := range defs {
		decoded := &SoftwareDef{}

		if val, ok := def["_app"]; ok {
			decoded.App = val.(string)
		}

		if _, ok := def["_deprecated"]; ok {
			decoded.Deprecated = true
		}

		if val, ok := def["_deps"]; ok {
			decoded.Deps = val.([]string)
		}

		if val, ok := def["_name"]; ok {
			decoded.Name = val.(string)
		}

		if val, ok := def["_service"]; ok {
			decoded.Service = val.(string)
		}

		if val, ok := def["_systemd"]; ok {
			decoded.Systemd = val.(string)
		}

		if val, ok := def["apk"]; ok {
		}

		// Complex, colon keys
		for k, v := range def {
			if k == "_bin" {
				decoded.Bin["_"] = v.(string)
			}

			if strings.HasPrefix(k, "_bin:") {
				s := strings.TrimPrefix(k, "_bin:")
				decoded.Bin[s] = v.(string)
			}

			if k == "_post" {
				decoded.Post["_"] = v.(string)
			}

			if strings.HasPrefix(k, "_post:") {
				s := strings.TrimPrefix(k, "_post:")
				decoded.Post[s] = v.(string)
			}

			if k == "_when" {
				decoded.When["_"] = v.(string)
			}

			if strings.HasPrefix(k, "_when:") {
				s := strings.TrimPrefix(k, "_when:")
				decoded.When[s] = v.(string)
			}
		}

		sd[pkg] = *decoded
	}

	return nil
}

type ChezmoiData struct {
	SoftwareGroups SoftwareGroups `toml:"softwareGroups"`
}

type SoftwareGroups map[string][]string

type ordersMsg []string

type recipesMsg SoftwarePackages

func getOrders(file string) tea.Cmd {
	return func() tea.Msg {
		fileData, fileErr := os.ReadFile(file)
		if fileErr != nil {
			return errMsg{fileErr}
		}

		var parsedYaml ChezmoiData

		yamlErr := yaml.Unmarshal(fileData, &parsedYaml)
		if yamlErr != nil {
			return errMsg{yamlErr}
		}

		return ordersMsg(parsedYaml.SoftwareGroups[softwareGroup])
	}
}

func getRecipes(file string) tea.Cmd {
	return func() tea.Msg {
		fileData, fileErr := os.ReadFile(file)
		if fileErr != nil {
			return errMsg{fileErr}
		}

		var parsedYaml SoftwarePackages

		yamlErr := yaml.Unmarshal(fileData, &parsedYaml)
		if yamlErr != nil {
			return errMsg{yamlErr}
		}

		return recipesMsg(parsedYaml)
	}
}
