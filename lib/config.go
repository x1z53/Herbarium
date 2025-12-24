package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"gopkg.in/yaml.v3"
)

const appConfigDir = "ru.ximper.Herbarium"

var (
	modsLineRe = regexp.MustCompile(
		`(?m)^\s*\$?\s*mods\s*\[\s*["']([^"'\]]+)["']\s*\]\s*=\s*u?["']([\s\S]*?)["']`,
	)
	braceRe = regexp.MustCompile(`\{[^}]*\}`)
)

func configDir() (string, error) {
	base, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(base, appConfigDir)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

func configPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.yaml"), nil
}

func modsDBPath() (string, error) {
	dir, err := configDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "mods_db.yaml"), nil
}

func EnsureConfig() (*Config, error) {
	if c, err := LoadConfig(); err == nil {
		return c, nil
	}

	home, _ := os.UserHomeDir()
	c := &Config{
		GameExe:     "/usr/bin/steam",
		Args:        []string{"-applaunch", "331470"},
		Root:        filepath.Join(home, ".steam/steam/steamapps/workshop/content/331470"),
		DisabledDir: filepath.Join(home, ".elmod_disabled"),
	}
	return c, SaveConfig(c)
}

func EnsureModsDB() (*ModsDB, error) {
	if db, err := LoadModsDB(); err == nil {
		return db, nil
	}

	db := &ModsDB{Mods: []ModEntry{}}
	return db, SaveModsDB(db)
}

func ToggleEnabled(enable bool, id string) error {
	if id == "" {
		return fmt.Errorf("provide folder id, codename, or ALL")
	}

	cfg, err := EnsureConfig()
	if err != nil {
		return err
	}

	db, err := EnsureModsDB()
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

	if err := SaveModsDB(db); err != nil {
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

func LoadConfig() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	return &c, yaml.Unmarshal(b, &c)
}

func SaveConfig(c *Config) error {
	path, err := configPath()
	if err != nil {
		return err
	}

	b, _ := yaml.Marshal(c)
	return os.WriteFile(path, b, 0644)
}

func LoadModsDB() (*ModsDB, error) {
	path, err := modsDBPath()
	if err != nil {
		return nil, err
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var db ModsDB
	return &db, yaml.Unmarshal(b, &db)
}

func SaveModsDB(db *ModsDB) error {
	b, _ := yaml.Marshal(db)
	path, err := modsDBPath()
	if err != nil {
		return err
	}

	return os.WriteFile(path, b, 0644)
}
