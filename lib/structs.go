package lib

import "time"

type Config struct {
	GameExe     string   `yaml:"game_exe"`
	Args        []string `yaml:"args,omitempty"`
	Root        string   `yaml:"workshop_root"`
	DisabledDir string   `yaml:"disabled_dir"`
}

type ModEntry struct {
	Name         string    `yaml:"name"`
	CodeName     string    `yaml:"codename"`
	Folder       string    `yaml:"folder"`
	Enabled      bool      `yaml:"enabled"`
	DiscoveredAt time.Time `yaml:"discovered_at"`
}

type ModsDB struct {
	Mods []ModEntry `yaml:"mods"`
}
