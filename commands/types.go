package commands

type Command = func(projectPath string, osArgs []string)
type CommandBuilderOutput struct {
	Cmd  Command
	Name string
}

type NitNode struct {
	Id    string
	Hash  string
	Files []*NitNode
	Type  string
}

type FileContent struct {
	Id   string
	Hash string
	Type string
}

type CommitObject struct {
	TreeHash string
	Parent   string
	Author   string
	Message  string
}
