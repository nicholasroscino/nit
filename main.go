package main

import (
	"log"
	"nit/commands"
	"nit/utils"
)

func main() {
	//rootFolder, err := os.Getwd()
	//utils.Check(err, "Unable to get current working directory\n")

	//nitFolder, err := utils.GetNitRepoFolder(rootFolder)

	//if err != nil {
	//	log.Fatal(err.Error())
	//}

	//val, err := commands.CatFileCommand(nitFolder, "3ccc5503fe7775a47905c5bf7be999189234d9c8")
	//
	//if err != nil {
	//	log.Fatal(err.Error())
	//
	//}

	//log.Println(val)
	//err = commands.AddCommand(nitFolder, "utils/utils.go")

	//aec33fcd5dc772b58e1235ccb90c17c204b13267 utils/utils.go 2025-05-17T22:19:30+02:00
	stagedObj := &utils.StagedObject{
		Hash:      "aec33fcd5dc772b58e1235ccb90c17c204b13267",
		Path:      "utils/culo/spectacle/prova/utils.go",
		Timestamp: "2025-05-17T22:19:30+02:00",
	}

	file := commands.WriteTree(stagedObj, nil)

	log.Println(file)

	// Check if the file was added to the index
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
}
