package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func prettyPrintAlias(alias Alias) {
	fmt.Println("trigger:")
	fmt.Printf("./tool %s\n", alias.Alias)
	fmt.Println("")
	fmt.Println("command:")
	fmt.Println(alias.Command)
	fmt.Println("")
	if len(alias.Path) > 0 {
		fmt.Println("executed in this path:")
		fmt.Println(alias.Path)
		fmt.Println("")
	}
}

func main() {
	rawJSON, err := ioutil.ReadFile("./config/aliases.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	aliasesBody := []byte(rawJSON)
	aliases := make([]Alias, 0)
	json.Unmarshal(aliasesBody, &aliases)

	list := flag.Bool("l", false, "List the aliases available")
	alias := flag.String("r", "", "Run alias")
	flag.Parse()

	if *list {
		fmt.Printf("Aliases:\n\n")
		for _, alias := range aliases {
			prettyPrintAlias(alias)
		}
	} else if len(*alias) > 0 {
		for _, aliasFromConfiguration := range aliases {
			if aliasFromConfiguration.Alias == *alias {
				command := exec.Command(aliasFromConfiguration.Command[0], aliasFromConfiguration.Command[1:]...)

				if len(aliasFromConfiguration.Path) > 0 {
					if aliasFromConfiguration.Path == "PWD" {
						command.Dir = os.Getenv("PWD")
					} else {
						command.Dir = aliasFromConfiguration.Path

					}
				}

				output, err := command.CombinedOutput()
				log.Printf("Running command and waiting for it to finish...")
				if err != nil {
					os.Stderr.WriteString(err.Error())
				}
				fmt.Println(string(output))
			}
		}
	} else {
		flag.PrintDefaults()
	}

	os.Exit(1)
}
