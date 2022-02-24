package cmd

import (
	"fmt"

	"github.com/go-clix/cli"
	"jensch.works/zl/cmd/zl/config"
)

func Switch() *cli.Command {
	c := &cli.Command{
		Use:   `switch`,
		Short: "Switch to another profile",
		Args:  cli.ArgsExact(1),
		Run: func(cmd *cli.Command, args []string) error {
			cfg, err := config.Default()
			if err != nil {
				return err
			}
			var match *config.Profile
			for i := range cfg.Profiles {
				if p := cfg.Profiles[i]; p.Name == args[0] {
					match = &p
					break
				}
			}
			if match == nil {
				return fmt.Errorf("profile not found: %q", args[0])
			}
			cfg.ActiveProfile = match.Name
			return config.SetDefault(cfg)
		},
	}
	return c
}
