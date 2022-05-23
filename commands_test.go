package main

import "testing"

func TestLoadCommands(t *testing.T) {
	cmds := loadCommands()
	t.Run(
		"there are 5 commands",
		func(t *testing.T) {
			if len(cmds) != 5 {
				t.Errorf("expected 5 commands, got %d commands", len(cmds))
			}
		},
	)
	t.Run(
		"1st command is indices",
		func(t *testing.T) {
			cmd := cmds[0]
			if cmd.Name != "indices" {
				t.Errorf("expected `indices` command, got `%s`", cmd.Name)
			}
			if len(cmd.Subcommands) != 4 {
				t.Errorf("expected 4 subcommands, got %d subcommands", len(cmd.Subcommands))
			}
		},
	)
	t.Run(
		"2nd command is cat",
		func(t *testing.T) {
			cmd := cmds[1]
			if cmd.Name != "cat" {
				t.Errorf("expected `cat` command, got `%s`", cmd.Name)
			}
			if len(cmd.Subcommands) != 2 {
				t.Errorf("expected 2 subcommands, got %d subcommands", len(cmd.Subcommands))
			}
		},
	)
	t.Run(
		"3rd command is cat",
		func(t *testing.T) {
			cmd := cmds[2]
			if cmd.Name != "search-template" {
				t.Errorf("expected `search-template` command, got `%s`", cmd.Name)
			}
			if len(cmd.Subcommands) != 5 {
				t.Errorf("expected 5 subcommands, got %d subcommands", len(cmd.Subcommands))
			}
		},
	)
	t.Run(
		"4th command is bulk",
		func(t *testing.T) {
			cmd := cmds[3]
			if cmd.Name != "bulk" {
				t.Errorf("expected `bulk` command, got `%s`", cmd.Name)
			}
			if len(cmd.Subcommands) != 1 {
				t.Errorf("expected 1 subcommand, got %d subcommands", len(cmd.Subcommands))
			}
		},
	)
	t.Run(
		"5th command is bulk",
		func(t *testing.T) {
			cmd := cmds[4]
			if cmd.Name != "reindex" {
				t.Errorf("expected `bulk` command, got `%s`", cmd.Name)
			}
			if len(cmd.Subcommands) != 0 {
				t.Errorf("expected no subcommands, got %d subcommands", len(cmd.Subcommands))
			}
		},
	)
}
