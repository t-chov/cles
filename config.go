package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/olivere/elastic/v7"
)

type profiles struct {
	Profile []EsConfig
}

type EsConfig struct {
	Name     string
	Address  []string
	Sniff    bool
	Username string
	Password string
}

func loadConfigFile(path string, profile string, debugFn func(msg string)) (*EsConfig, error) {
	debugFn(fmt.Sprintf("config file path: %s", path))
	f, err := ioutil.ReadFile(path)
	if err == nil {
		profiles := &profiles{}
		err := toml.Unmarshal(f, profiles)
		if err != nil {
			return nil, err
		}
		for _, prof := range profiles.Profile {
			if prof.Name == profile {
				cfg := EsConfig{
					Name:     prof.Name,
					Address:  prof.Address,
					Username: prof.Username,
					Password: prof.Password,
					Sniff:    prof.Sniff,
				}
				debugOut, _ := json.MarshalIndent(cfg, "", "  ")
				debugFn(fmt.Sprintf("config: %s", debugOut))
				return &cfg, nil
			}
		}
	}
	return nil, nil
}

func (cfg *EsConfig) Load(profileName string, debugFn func(message string)) error {
	// default value
	cfg.Address = []string{"http://127.0.0.1:9200"}
	cfg.Username = ""
	cfg.Password = ""
	cfg.Sniff = false

	// load from config file.
	// windows: $APPDATA\cles\
	// unix-like: $HOME/.config/cles/
	var dir string
	if runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	cfgPath := filepath.Join(dir, "cles", "config.toml")
	debugFn(fmt.Sprintf("load config from %s", cfgPath))

	f, err := ioutil.ReadFile(cfgPath)
	if err == nil {
		profiles := &profiles{}
		err := toml.Unmarshal(f, profiles)
		if err != nil {
			return err
		}
		for _, profile := range profiles.Profile {
			if profile.Name == profileName {
				cfg.Name = profile.Name
				cfg.Address = profile.Address
				cfg.Username = profile.Username
				cfg.Password = profile.Password
				cfg.Sniff = profile.Sniff
				debugOut, _ := json.MarshalIndent(cfg, "", "  ")
				debugFn(fmt.Sprintf("config: %s", debugOut))
				return nil
			}
		}
	}

	// overwrite with environment variables
	if len(os.Getenv("ES_ADDRESS")) > 0 {
		cfg.Address = strings.Split(os.Getenv("ES_ADDRESS"), ",")
	}
	if len(os.Getenv("ES_USERNAME")) > 0 {
		cfg.Username = os.Getenv("ES_USERNAME")
	}
	if len(os.Getenv("ES_PASSWORD")) > 0 {
		cfg.Password = os.Getenv("ES_PASSWORD")
	}
	if len(os.Getenv("ES_SNIFF")) > 0 {
		cfg.Sniff = strings.ToLower(os.Getenv("ES_SNIFF")) == "true"
	}
	debugOut, _ := json.MarshalIndent(cfg, "", "  ")
	debugFn(fmt.Sprintf("config: %s", debugOut))

	return nil
}

func initClient(profile string, debugFn func(message string)) (*elastic.Client, error) {
	var cfg EsConfig
	err := cfg.Load(profile, debugFn)
	if err != nil {
		return nil, err
	}
	debugFn(fmt.Sprintf("URL  : %v", cfg.Address))
	debugFn(fmt.Sprintf("USER : %s", cfg.Username))
	debugFn(fmt.Sprintf("PASS : %s", cfg.Password))
	debugFn(fmt.Sprintf("SNIFF: %v", cfg.Sniff))
	client, err := elastic.NewClient(
		elastic.SetURL(cfg.Address...),
		elastic.SetBasicAuth(cfg.Username, cfg.Password),
		elastic.SetSniff(cfg.Sniff),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
