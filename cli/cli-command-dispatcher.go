package cli

import (
	"log"
	"nit/commands"
	"nit/commands/add"
	"nit/commands/cat"
	"nit/commands/commit"
	"nit/commands/hash-object"
	init_pkg "nit/commands/init"
	nitlog "nit/commands/log"
	"nit/commands/write-tree"
)

type CliCommandDispatcher struct {
	commands map[string]commands.Command
}

func NewCliCommandDispatcher() *CliCommandDispatcher {
	return &CliCommandDispatcher{
		commands: make(map[string]commands.Command),
	}
}

func (r *CliCommandDispatcher) Init() {
	hashObjectHandler := hash_object.CommandBuilder()
	commitHandler := commit.CommandBuilder()
	initHandler := init_pkg.CommandBuilder()
	addHandler := add.CommandBuilder()
	catHandler := cat.CommandBuilder()
	writeTreeHandler := write_tree.CommandBuilder()
	logHandler := nitlog.CommandBuilder()

	r.commands[hashObjectHandler.Name] = hashObjectHandler.Cmd
	r.commands[commitHandler.Name] = commitHandler.Cmd
	r.commands[initHandler.Name] = initHandler.Cmd
	r.commands[addHandler.Name] = addHandler.Cmd
	r.commands[catHandler.Name] = catHandler.Cmd
	r.commands[writeTreeHandler.Name] = writeTreeHandler.Cmd
	r.commands[logHandler.Name] = logHandler.Cmd
}

func (r *CliCommandDispatcher) DispatchCommand(projectPath string, osArgs []string) {
	command := r.commands[osArgs[1]]

	if command == nil {
		log.Fatal("Command not found: " + osArgs[1] + "\n")
	}

	command(projectPath, osArgs[2:])
}
