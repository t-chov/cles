package main

import (
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

func setConfigFilePath(runtimeOS string) string {
	// set config file by OS.
	// windows: $APPDATA\cles\
	// unix-like: $HOME/.config/cles/
	var dir string
	if runtimeOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data")
		}
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	return filepath.Join(dir, "cles", "config.toml")
}

func loadConfigFromFile(path string, profile string) (*EsConfig, error) {
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
				return &cfg, nil
			}
		}
	}
	return nil, nil
}

func loadConfig(cfgPath string, profileName string) (*EsConfig, error) {
	// default value
	cfg := EsConfig{
		Name:     profileName,
		Address:  []string{"http://127.0.0.1:9200"},
		Username: "",
		Password: "",
		Sniff:    false,
	}

	fcfg, err := loadConfigFromFile(cfgPath, profileName)
	if err != nil {
		return nil, err
	}
	if fcfg != nil {
		cfg = *fcfg
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

	return &cfg, nil
}

func initClient(profile string, debugFn func(message string)) (*elastic.Client, error) {
	cfgPath := setConfigFilePath(runtime.GOOS)
	debugFn(fmt.Sprintf("config file path: %s", cfgPath))
	cfg, err := loadConfig(cfgPath, profile)
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
