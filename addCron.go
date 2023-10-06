package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func main() {
	cronName := flag.String("cronName", "", "Name of the cron file em /etc/cron.d/")
	cron := flag.String("cron", "", "Cron expression, example 0 9,12,15,18,21 * * * root /bin/my-script")
	description := flag.String("description", "", "Description of the cron job (optional)")
	help := flag.Bool("help", false, "Show the available args params")

	flag.Parse()

	if !isArgsOk(*help) {
		return
	}

	validate(cronName, cron, description)
}

func validate(cronName *string, cron *string, description *string) {
	if *cronName == "" || *cron == "" {
		isArgsOk(true)
		os.Exit(1)
	}

	fileName := fmt.Sprintf("/etc/cron.d/%s", *cronName)

	fileContents := checkFile(fileName)

	if strings.Contains(string(fileContents), *cron) {
		fmt.Printf("Cron already on %s\n", fileName)
		readFileContent(fileName)
	} else {
		writeCronOnFile(description, cron, fileName, fileContents)
	}
}

func writeCronOnFile(description *string, cron *string, fileName string, fileContents []byte) {
	comment := fmt.Sprintf("\n# %s\n%s\n", *description, *cron)
	err := ioutil.WriteFile(fileName, []byte(string(fileContents)+comment), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("Cron registered")
	readFileContent(fileName)
}

func checkFile(fileName string) []byte {
	tryCreateFile(fileName)

	fileContents, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		os.Exit(1)
	}
	return fileContents
}

func isArgsOk(help bool) bool {
	if help {
		fmt.Printf("Available args params...\n\n")
		fmt.Printf("\t%s: %s\n", "--cronName", "The file that with registered the cron. If file not exists, will be created")
		fmt.Printf("\t%s: %s\n", "--cron", "Cron expression, example 0 9,12,15,18,21 * * * root /bin/my-script")
		fmt.Printf("\t%s: %s\n\n", "--description", "Description of the cron job (optional)")
		return false
	}
	return true
}

func tryCreateFile(fileName string) {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Printf("Error creating file: %s\n", err)
			os.Exit(1)
		}
		_ = file.Close()
	}
}

func readFileContent(fileName string) {
	newFileContents, _ := ioutil.ReadFile(fileName)
	fmt.Println(string(newFileContents))
}
