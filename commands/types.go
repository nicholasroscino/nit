package commands

type Command = func(projectPath string, osArgs []string)
type CommandBuilderOutput struct {
	Cmd  Command
	Name string
}
