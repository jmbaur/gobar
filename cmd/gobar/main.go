package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

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
		case "memory":
			mod = &module.Memory{}
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

func must(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	exe, err := os.Executable()
	must(err)

	log.SetPrefix(fmt.Sprintf("%s: ", filepath.Base(exe)))
	log.SetFlags(log.Lmsgprefix)
	log.Println("running")

	configFile := flag.String("config", "", "Path to gobar.yaml config file")
	flag.Parse()

	cfg, err := config.GetConfig(*configFile)
	must(err)

	must(module.Run(decodeToModules(cfg)...))
}
