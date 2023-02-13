package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const mainPathConst = "{{MainPath}}"

type program struct {
	Command string
	Args    []string
}

type commonConfig struct {
	Key   string
	Value string
}

type task struct {
	Name         string
	MainPath     string
	Programs     []program
	CommonConfig []commonConfig
}

type toExecute struct {
	Tasks []task
}

func main() {
	//get the config file location from the args
	configFile := os.Args[2]
	if configFile == "" {
		fmt.Println("Error: No config file provided")
	}

	// loads the config file and unmarshals it into the toExecute struct
	file, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Println("Error: ", err)
	}
	// unmarshal the file into the toExecute struct
	var toExecute toExecute
	if err := json.Unmarshal(file, &toExecute); err != nil {
		fmt.Println("Error: ", err)
	}

	for _, task := range toExecute.Tasks {
		// join the main path and the config file using the filepath package
		generatedCommonFile := filepath.Join(task.MainPath, "Config", "common.ini")

		file, err := ioutil.ReadFile("./mt5-common-file-template.ini")
		if err != nil {
			fmt.Println("Error: ", err)
		}
		common := string(file)
		for _, config := range task.CommonConfig {
			key := fmt.Sprintf("{{%s}}", config.Key)
			common = strings.ReplaceAll(common, key, config.Value)
		}
		if err := ioutil.WriteFile(generatedCommonFile, []byte(common), 0644); err != nil {
			fmt.Println("Error: ", err)
		}
		for _, program := range task.Programs {
			// execute the program
			fmt.Printf("Running program: %s", program.Command)
			args := []string{"/C", strings.ReplaceAll(program.Command, mainPathConst, task.MainPath)}
			for _, arg := range program.Args {
				args = append(args, strings.ReplaceAll(arg, mainPathConst, task.MainPath))
			}
			c := exec.Command("cmd", args...)
			if err := c.Run(); err != nil {
				fmt.Println("Error: ", err)
			}
			c = exec.Command("cmd", "/C", "del", "D:\\a.txt")

			if err := c.Run(); err != nil {
				fmt.Println("Error: ", err)
			}
		}
	}
}
