package cmd

import (
	"flag"

	"github.com/jmbaur/gobar/config"
	"github.com/jmbaur/gobar/module"
	"github.com/mitchellh/mapstructure"
)

func decodeToModules(cfg *config.Config) []module.Module {
	modules := []module.Module{}

	for _, m := range cfg.Modules {
		maybeMod, ok := m.(map[any]any)
		if !ok {
			continue
		}

		name, ok := maybeMod["module"].(string)
		if !ok {
			continue
		}

		var mod module.Module
		switch name {
		case "battery":
			mod = &module.Battery{}
		case "datetime":
			mod = &module.Datetime{}
		case "network":
			mod = &module.Network{}
		case "text":
			mod = &module.Text{}
		default:
			continue
		}

		if err := mapstructure.Decode(m, &mod); err != nil {
			continue
		}
		modules = append(modules, mod)
	}

	return modules
}

func Run() error {
	configFile := flag.String("config", "", "Path to gobar.yaml config file")
	flag.Parse()

	cfg, err := config.GetConfig(*configFile)
	if err != nil {
		return err
	}

	return module.Run(decodeToModules(cfg)...)
}
