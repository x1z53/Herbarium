package main

import (
	"context"
	"fmt"
	"os"

	"esmodmanager/lib"

	"github.com/urfave/cli/v3"
)

func main() {
	const (
		configPath = "config.yaml"
		dbPath     = "mods_db.yaml"
	)

	cmd := &cli.Command{
		Name:  "esmodmanager",
		Usage: "Manager for Everlasting Summer mods",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List known mods",
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg, err := lib.EnsureConfig(configPath)
					if err != nil {
						return err
					}

					db, err := lib.EnsureModsDB(dbPath)
					if err != nil {
						return err
					}

					lib.ScanAndUpdate(cfg, db)

					if err := lib.SaveModsDB(dbPath, db); err != nil {
						return err
					}

					lib.PrintMods(db)
					return nil
				},
			},

			{
				Name:      "disable",
				Aliases:   []string{"d"},
				Usage:     "Disable mod by numeric folder, codename, or ALL",
				ArgsUsage: "<id>",
				Action: func(ctx context.Context, c *cli.Command) error {
					return lib.ToggleEnabled(configPath, dbPath, false, c.Args().First())
				},
			},

			{
				Name:      "enable",
				Aliases:   []string{"e"},
				Usage:     "Enable mod by numeric folder, codename, or ALL",
				ArgsUsage: "<id>",
				Action: func(ctx context.Context, c *cli.Command) error {
					return lib.ToggleEnabled(configPath, dbPath, true, c.Args().First())
				},
			},

			{
				Name:    "launch",
				Aliases: []string{"start", "l"},
				Usage:   "Launch game with current mod setup",
				Action: func(ctx context.Context, c *cli.Command) error {
					cfg, err := lib.EnsureConfig(configPath)
					if err != nil {
						return err
					}

					if cfg.GameExe == "" {
						return fmt.Errorf("game_exe is empty in %s", configPath)
					}

					db, err := lib.EnsureModsDB(dbPath)
					if err != nil {
						return err
					}

					lib.ScanAndUpdate(cfg, db)

					if err := lib.SaveModsDB(dbPath, db); err != nil {
						return err
					}

					lib.LaunchWithMods(cfg, db)
					return nil
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
