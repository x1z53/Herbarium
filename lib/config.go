package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

var (
	modsLineRe = regexp.MustCompile(
		`(?m)^\s*\$?\s*mods\s*\[\s*["']([^"'\]]+)["']\s*\]\s*=\s*u?["']([\s\S]*?)["']`,
	)
	braceRe = regexp.MustCompile(`\{[^}]*\}`)
)

func EnsureConfig(path string) (*Config, error) {
	if c, err := LoadConfig(path); err == nil {
		return c, nil
	}

	home, _ := os.UserHomeDir()
	c := &Config{
		GameExe:     "/usr/bin/steam",
		Args:        []string{"-applaunch", "331470"},
		Root:        filepath.Join(home, ".steam/steam/steamapps/workshop/content/331470"),
		DisabledDir: filepath.Join(home, ".elmod_disabled"),
	}
	return c, SaveConfig(path, c)
}

func EnsureModsDB(path string) (*ModsDB, error) {
	if db, err := LoadModsDB(path); err == nil {
		return db, nil
	}
	db := &ModsDB{Mods: []ModEntry{}}
	return db, SaveModsDB(path, db)
}

func ToggleEnabled(cfgPath, dbPath string, enable bool, id string) error {
	if id == "" {
		return fmt.Errorf("provide folder id, codename, or ALL")
	}

	cfg, err := EnsureConfig(cfgPath)
	if err != nil {
		return err
	}

	db, err := EnsureModsDB(dbPath)
	if err != nil {
		return err
	}

	ScanAndUpdate(cfg, db)

	if id == "ALL" {
		setAllEnabled(db, enable)
	} else {
		if err := setEnabled(db, id, enable); err != nil {
			return err
		}
	}

	if err := SaveModsDB(dbPath, db); err != nil {
		return err
	}

	action := map[bool]string{true: "Enabled", false: "Disabled"}[enable]
	fmt.Printf("%s %s\n", action, id)
	return nil
}

func setEnabled(db *ModsDB, id string, enabled bool) error {
	for i, m := range db.Mods {
		if m.Folder == id || m.CodeName == id {
			db.Mods[i].Enabled = enabled
			return nil
		}
	}
	return fmt.Errorf("mod not found: %s", id)
}

func setAllEnabled(db *ModsDB, enabled bool) {
	for i := range db.Mods {
		db.Mods[i].Enabled = enabled
	}
}

func LoadConfig(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	return &c, yaml.Unmarshal(b, &c)
}

func SaveConfig(path string, c *Config) error {
	b, _ := yaml.Marshal(c)
	return os.WriteFile(path, b, 0644)
}

func LoadModsDB(path string) (*ModsDB, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var db ModsDB
	return &db, yaml.Unmarshal(b, &db)
}

func SaveModsDB(path string, db *ModsDB) error {
	b, _ := yaml.Marshal(db)
	return os.WriteFile(path, b, 0644)
}
