package cmd

import (
	"flag"

	"github.com/jmbaur/gobar/config"
	"github.com/jmbaur/gobar/module"
	"github.com/mitchellh/mapstructure"
)

func Run() error {
	configFile := flag.String("config", "", "Path to gobar.yaml config file")
	flag.Parse()

	cfg, err := config.GetConfig(*configFile)
	if err != nil {
		return err
	}

	module.Run(decodeToModules(cfg)...)

	return nil
}

func decodeToModules(cfg *config.Config) []module.Module {
	modules := []module.Module{}

	for _, m := range cfg.Modules {
		casted, ok := m.(map[interface{}]interface{})
		if !ok {
			continue
		}

		moduleType, ok := casted["module"]
		if !ok {
			continue
		}

		// TODO(jared): use generics
		switch moduleType {
		case "network":
			var mod module.Network
			if err := mapstructure.Decode(m, &mod); err != nil {
				continue
			}
			modules = append(modules, mod)
		case "datetime":
			var mod module.Datetime
			if err := mapstructure.Decode(m, &mod); err != nil {
				continue
			}
			modules = append(modules, mod)
		case "battery":
			var mod module.Battery
			if err := mapstructure.Decode(m, &mod); err != nil {
				continue
			}
			modules = append(modules, mod)
		}
	}

	return modules
}
