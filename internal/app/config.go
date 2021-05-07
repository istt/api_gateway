package app

import (
	"log"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Config hold global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var Config = koanf.New(".")

// ConfigINit configure application runtime
func ConfigInit(path string) {
	// override configuration with YAML
	err := Config.Load(file.Provider(path), yaml.Parser())
	if err != nil {
		log.Printf("KoanfInit error config load: %s", err)
	} else {
		log.Printf("KoanfInit config load ok: %s", path)
	}
}

// ConfigWatch let the application watch for changes
func ConfigWatch(path string) {
	file.Provider(path).Watch(func(event interface{}, err error) {
		if err != nil {
			log.Printf("watch error: %v", err)
			return
		}

		log.Printf("config changed. Reloading %+v", event)
		Config.Load(file.Provider(path), yaml.Parser())
		Config.Print()
	})
}
