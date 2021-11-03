package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

func main() {
	outputDirRead, _ := os.Open("./../project")
	outputDirFiles, _ := outputDirRead.Readdir(0)

	// Loop over files.
	for _, val := range outputDirFiles {
		if !val.IsDir() {
			continue
		}

		path := "./../project/" + val.Name()

		outputDirReadInner, _ := os.Open(path)
		outputDirFilesInner, _ := outputDirReadInner.Readdir(0)

		var hasInner bool
		for _, innerPath := range outputDirFilesInner {
			if !innerPath.IsDir() {
				break
			}

			hasInner = true

			from := path + "/" + innerPath.Name()
			to := path + innerPath.Name()
			_, err := exec.Command("mv", from, to).Output()
			if err != nil {
				fmt.Printf("error %s", err)
			}
		}

		if hasInner {
			_, err := exec.Command("rm -f", path).Output()
			if err != nil {
				fmt.Printf("error %s", err)
			}
		}

		log.Println("______")
	}
}
