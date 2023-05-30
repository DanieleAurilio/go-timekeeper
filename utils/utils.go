package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"timekeeper/models"
)

const separator string = "/"
const ext string = ".go"

func ExecuteJobCmd(fullFilePath string, function string, params map[string]string) {
	outputDir := ReadConfigFile().OutputDir
	cmdExec := exec.Command("go", "build","-o", outputDir, fullFilePath)
	
	err := cmdExec.Run()

	if err != nil {
		fmt.Println("Error building job file " + fullFilePath + " - Error: " + err.Error())
		return
	}

	filename := fullFilePath[strings.LastIndex(fullFilePath, separator):strings.LastIndex(fullFilePath, ".")]
	outputFilepath := getOutputFilePath(filename) 

	var formatParams string
	if len(params) > 0 {
		jsonParamsByte, err := json.Marshal(params)

		if err != nil {
			fmt.Println(err.Error())
			return
		}

		formatParams = string(jsonParamsByte)
	}
	
	

	cmdRun := exec.Command(outputFilepath, function, formatParams)

	cmdRun.Stdout = os.Stdout
	cmdRun.Stderr = os.Stderr
	fmt.Println("Executing file " + outputFilepath)

	errRun := cmdRun.Run()

	if errRun != nil {
		fmt.Println("Error running job file " + outputFilepath + " Error: " + errRun.Error())
		return
	}

}

func ReadConfigFile() *models.Config {
	
	var job = &models.Config{}


	config, err := os.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	json.Unmarshal(config, &job)
	
	return job	
}

func ResolveJobPath(filename string) string {
	if !strings.Contains(filename, ext) {
		fmt.Println("Filename with " + ext + " ext not found")
		return ""
	}

	var fullFilePath string
	directoryPath, _ := os.Getwd()
	err := filepath.Walk(directoryPath, func (path string, info fs.FileInfo, err error) error  {
		if strings.Contains(path, filename) {
			fullFilePath = path
		}
		return nil
	})

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	if len(fullFilePath) == 0 {
		fmt.Println("No file with filename " + filename + " found.")
		return ""
	}
	
	return fullFilePath
}

func Mkdir(outputDirectoryPath string) {

	fileInfo, _ := os.Stat(outputDirectoryPath)
	
	if fileInfo != nil {
		return
	}
	
	err := os.MkdirAll(outputDirectoryPath, fs.ModeAppend)

	if err != nil {
		fmt.Println("Error during creating output directory: " + outputDirectoryPath + " Error: " + err.Error() )
	} else {
		fmt.Println("Output Directory : " + outputDirectoryPath + " Correctly created")
	}
}

func CheckTime(timeNow time.Time, jobTime time.Time) bool {
	return timeNow.Equal(jobTime)
}

func BuildOutputDirectoryPath(directoryName string) (error, string) {
	if len(directoryName) == 0 {
		fmt.Println("please provide a directory name in config.json file")
		return errors.New("please provide a directory name in config.json file"), ""
	}

	cwd, _ := os.Getwd()
	elems := make([]string, 0)
	elems = append(elems, cwd, directoryName)
	outputDirectoryPath := strings.Join(elems, separator)
	
	return nil, outputDirectoryPath
}

func getOutputFilePath(filename string) string {
	cwd, _ := os.Getwd()
	outputDir := ReadConfigFile().OutputDir
	elems := make([]string, 0)
	elems = append(elems, cwd, outputDir, filename)
	return  strings.Join(elems, separator)
}