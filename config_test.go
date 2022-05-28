package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	t.Run(
		"file note exist",
		func(t *testing.T) {
			cfg, err := loadConfigFromFile("/path/to/invalid", "notfound")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil, got %v", *cfg)
			}
		},
	)
	t.Run(
		"unmarshal error",
		func(t *testing.T) {
			cwd, _ := os.Getwd()
			path := filepath.Join(cwd, "main.go")
			cfg, err := loadConfigFromFile(path, "invalid")
			if err == nil {
				t.Errorf("expected error, got nil")
			}
			if !strings.Contains(fmt.Sprintf("%s", err), "toml") {
				t.Errorf("expected toml error, got %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil, got %v", *cfg)
			}
		},
	)
	t.Run(
		"matched profile",
		func(t *testing.T) {
			cwd, _ := os.Getwd()
			path := filepath.Join(cwd, ".config", "cles", "config.toml")
			cfg, err := loadConfigFromFile(path, "default")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if cfg == nil {
				t.Errorf("expected config, got nil")
				return
			}
			expected := EsConfig{
				Name:     "default",
				Address:  []string{"http://localhost:9200"},
				Sniff:    false,
				Username: "elastic",
				Password: "elasticPassword",
			}
			if !reflect.DeepEqual(expected, *cfg) {
				t.Errorf("expected %v, got %v", expected, *cfg)
			}
		},
	)
	t.Run(
		"unmatched profile",
		func(t *testing.T) {
			cwd, _ := os.Getwd()
			path := filepath.Join(cwd, ".config", "cles", "config.toml")
			cfg, err := loadConfigFromFile(path, "noprofile")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil config, got %v", cfg)
			}
		},
	)
}

func TestSetConfigFilePath(t *testing.T) {
	home := os.Getenv("HOME")
	userProfile := os.Getenv("USERPROFILE")
	appData := os.Getenv("APPDATA")
	t.Run(
		"unix-like",
		func(t *testing.T) {
			err := os.Setenv("HOME", "/foo/bar")
			if err != nil {
				t.Errorf("os.Setenv error: %v", err)
			}
			got := setConfigFilePath("ubuntu")
			expected := "/foo/bar/.config/cles/config.toml"
			if expected != got {
				t.Errorf("expected %s, got %s", expected, got)
			}
			os.Setenv("HOME", home)
		},
	)
	t.Run(
		"windows, no APPDATA",
		func(t *testing.T) {
			err := os.Setenv("USERPROFILE", "/user/profile")
			if err != nil {
				t.Errorf("os.Setenv error: %v", err)
			}
			got := setConfigFilePath("windows")
			expected := "/user/profile/Application Data/cles/config.toml"
			if expected != got {
				t.Errorf("expected %s, got %s", expected, got)
			}
			os.Setenv("USERPROFILE", userProfile)
		},
	)
	t.Run(
		"windows, with APPDATA",
		func(t *testing.T) {
			err := os.Setenv("APPDATA", "/appdata")
			if err != nil {
				t.Errorf("os.Setenv error: %v", err)
			}
			got := setConfigFilePath("windows")
			expected := "/appdata/cles/config.toml"
			if expected != got {
				t.Errorf("expected %s, got %s", expected, got)
			}
			os.Setenv("APPDATA", appData)
		},
	)
}

func TestLoadConfig(t *testing.T) {
	t.Run(
		"no config, no envvar",
		func(t *testing.T) {
			cfg, err := loadConfig("/invalid/path", "noprofile")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			expected := EsConfig{
				Name:     "noprofile",
				Address:  []string{"http://127.0.0.1:9200"},
				Username: "",
				Password: "",
				Sniff:    false,
			}
			if !reflect.DeepEqual(expected, *cfg) {
				t.Errorf("expected %v, got %v", expected, *cfg)
			}
		},
	)
	t.Run(
		"successful load from config file",
		func(t *testing.T) {
			cwd, err := os.Getwd()
			if err != nil {
				t.Errorf("os.Getwd error: %v", err)
				return
			}
			cfgPath := filepath.Join(cwd, ".config", "cles", "config.toml")
			cfg, err := loadConfig(cfgPath, "default")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}
			expected := EsConfig{
				Name:     "default",
				Address:  []string{"http://localhost:9200"},
				Username: "elastic",
				Password: "elasticPassword",
				Sniff:    false,
			}
			if !reflect.DeepEqual(expected, *cfg) {
				t.Errorf("expected %v, got %v", expected, *cfg)
			}
		},
	)
	t.Run(
		"loadconfig error",
		func(t *testing.T) {
			cwd, err := os.Getwd()
			if err != nil {
				t.Errorf("os.Getwd error: %v", err)
				return
			}
			cfgPath := filepath.Join(cwd, "README.md")
			cfg, err := loadConfig(cfgPath, "default")
			if err == nil {
				t.Errorf("expected error, got nil")
				return
			}
			if !strings.Contains(fmt.Sprintf("%s", err), "toml") {
				t.Errorf("expected toml error, got %v", err)
			}
			if cfg != nil {
				t.Errorf("expected nil config, got %v", *cfg)
			}
		},
	)
	t.Run(
		"successful load from config file and overwrite by envvar",
		func(t *testing.T) {
			originalEnvs := map[string]string{
				"ES_ADDRESS":  os.Getenv("ES_ADDRESS"),
				"ES_USERNAME": os.Getenv("ES_USERNAME"),
				"ES_PASSWORD": os.Getenv("ES_PASSWORD"),
				"ES_SNIFF":    os.Getenv("ES_SNIFF"),
			}
			testEnvs := map[string]string{
				"ES_ADDRESS":  "http://example00.es,http://example01.es",
				"ES_USERNAME": "elasticUser",
				"ES_PASSWORD": "elasticPassword00",
				"ES_SNIFF":    "true",
			}
			cwd, err := os.Getwd()
			if err != nil {
				t.Errorf("os.Getwd error: %v", err)
				return
			}
			cfgPath := filepath.Join(cwd, ".config", "cles", "config.toml")
			for key, value := range testEnvs {
				err := os.Setenv(key, value)
				if err != nil {
					t.Errorf("os.Setenv error: %v", err)
					return
				}
			}
			cfg, err := loadConfig(cfgPath, "default")
			if err != nil {
				t.Errorf("expected no error, got %v", err)
				return
			}
			expected := EsConfig{
				Name:     "default",
				Address:  []string{"http://example00.es", "http://example01.es"},
				Username: "elasticUser",
				Password: "elasticPassword00",
				Sniff:    true,
			}
			if !reflect.DeepEqual(expected, *cfg) {
				t.Errorf("expected %v, got %v", expected, *cfg)
			}
			for key, value := range originalEnvs {
				os.Setenv(key, value)
			}
		},
	)
}
