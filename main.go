package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const baseDir = "./project/"
const projectImportName = `"autokz-register/`

func main() {
	outputDirRead, _ := os.Open(baseDir)
	outputDirFiles, _ := outputDirRead.Readdir(0)

	// Loop over files.
	for _, val := range outputDirFiles {
		if !val.IsDir() {
			continue
		}

		path := baseDir + val.Name()

		outputDirReadInner, _ := os.Open(path)
		outputDirFilesInner, _ := outputDirReadInner.Readdir(0)

		var hasInner bool
		for _, innerPath := range outputDirFilesInner {
			if !innerPath.IsDir() {
				break
			}

			hasInner = true

			from := path + "/" + innerPath.Name()
			to := path + toCamelInitCase(innerPath.Name(), true)
			_, err := exec.Command("mv", from, to).Output()
			if err != nil {
				fmt.Printf("error %s", err)
			}

			// Fix imports...
			_ = filepath.Walk(to, func(path string, info os.FileInfo, _ error) error {
				if info.IsDir() {
					return nil
				}

				// Read
				input, err := ioutil.ReadFile(path)
				if err != nil {
					log.Fatalln(err)
				}

				// Write
				oldImport := projectImportName + val.Name() + "/" + innerPath.Name()
				newImport := projectImportName + val.Name() + toCamelInitCase(innerPath.Name(), true)

				output := []byte(strings.ReplaceAll(string(input), oldImport, newImport))
				if err = ioutil.WriteFile(path, output, 0644); err != nil {
					log.Fatalln(err)
				}

				return nil
			})
		}

		if hasInner {
			if err := os.RemoveAll(path); err != nil {
				log.Fatalf("Error remove dir: " + err.Error())
			}
		}
	}
}

// Converts a string to CamelCase
func toCamelInitCase(s string, initCase bool) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	n := strings.Builder{}
	n.Grow(len(s))
	capNext := initCase
	for i, v := range []byte(s) {
		vIsCap := v >= 'A' && v <= 'Z'
		vIsLow := v >= 'a' && v <= 'z'
		if capNext {
			if vIsLow {
				v += 'A'
				v -= 'a'
			}
		} else if i == 0 {
			if vIsCap {
				v += 'a'
				v -= 'A'
			}
		}
		if vIsCap || vIsLow {
			n.WriteByte(v)
			capNext = false
		} else if vIsNum := v >= '0' && v <= '9'; vIsNum {
			n.WriteByte(v)
			capNext = true
		} else {
			capNext = v == '_' || v == ' ' || v == '-' || v == '.'
		}
	}
	return n.String()
}
