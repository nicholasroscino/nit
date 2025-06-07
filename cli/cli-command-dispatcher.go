package cli

import (
	"log"
	"nit/commands"
	"nit/commands/add"
	"nit/commands/cat"
	"nit/commands/commit"
	"nit/commands/hash-object"
	initpkg "nit/commands/init"
	nitlog "nit/commands/log"
	"nit/commands/walk"
	"nit/commands/write-tree"
)

type CommandDispatcher struct {
	commands map[string]commands.Command
}

func NewCommandDispatcher() *CommandDispatcher {
	return &CommandDispatcher{
		commands: make(map[string]commands.Command),
	}
}

func (r *CommandDispatcher) Init() {
	hashObjectHandler := hash_object.CommandBuilder()
	commitHandler := commit.CommandBuilder()
	initHandler := initpkg.CommandBuilder()
	addHandler := add.CommandBuilder()
	catHandler := cat.CommandBuilder()
	writeTreeHandler := write_tree.CommandBuilder()
	logHandler := nitlog.CommandBuilder()
	walkHandler := walk.CommandBuilder()

	r.commands[hashObjectHandler.Name] = hashObjectHandler.Cmd
	r.commands[commitHandler.Name] = commitHandler.Cmd
	r.commands[initHandler.Name] = initHandler.Cmd
	r.commands[addHandler.Name] = addHandler.Cmd
	r.commands[catHandler.Name] = catHandler.Cmd
	r.commands[writeTreeHandler.Name] = writeTreeHandler.Cmd
	r.commands[logHandler.Name] = logHandler.Cmd
	r.commands[walkHandler.Name] = walkHandler.Cmd
}

func (r *CommandDispatcher) DispatchCommand(projectPath string, osArgs []string) {
	command := r.commands[osArgs[1]]

	if command == nil {
		log.Fatal("Command not found: " + osArgs[1] + "\n")
	}

	command(projectPath, osArgs[2:])
}
