package main

import "testing"

func TestInitApp(t *testing.T) {
	t.Run(
		"initiate application",
		func(t *testing.T) {
			app := initApp()
			if app.Name != "cles" {
				t.Errorf("expect `cles`, got `%s`", app.Name)
			}
		},
	)
}
