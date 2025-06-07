package commands

import "strings"

type Command = func(projectPath string, osArgs []string)
type CommandBuilderOutput struct {
	Cmd  Command
	Name string
}

type NitTreeFileContent struct {
	Type     string
	Hash     string
	PathName string
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

func (c *CommitObject) BeautyPrint() string {
	stringa := "Tree: " + c.TreeHash + "\n"

	if c.Parent != "" {
		stringa += "Parent: " + c.Parent + "\n"
	} else {
		stringa += "Parent: (none)\n"
	}

	stringa += "Author: " + c.Author + "\n" +
		"Message: " + c.Message + "\n"

	return stringa
}

func NewCommitObject(fileContent string) *CommitObject {
	lines := strings.Split(fileContent, "\n")

	var values = make(map[string]string)
	for _, line := range lines {
		lineVals := strings.Split(line, " ")

		if len(lineVals) < 2 {
			continue
		}

		key := strings.TrimSpace(lineVals[0])
		value := strings.TrimSpace(strings.Join(lineVals[1:], " "))
		values[key] = value
	}

	commitMessage := strings.TrimSpace(lines[len(lines)-1]) // Trim the last line

	return &CommitObject{
		TreeHash: values["tree"],
		Parent:   values["parent"],
		Author:   values["author"],
		Message:  commitMessage,
	}
}
