package service

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sysu-go-online/service-end/model/entities"
)

// ROOT defines the root directory
var ROOT = "/home/golang/src/github.com"

// GetProjectNameByID return project according to the given id
func GetProjectNameByID(id string) (string, error) {
	return "test", nil
}

// UpdateFileContent update content with given filepath and content
func UpdateFileContent(projectid string, filePath string, content string) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		panic(err)
	}
	absPath := filepath.Join(ROOT, projectName, filePath)

	// Update file, if the file not exists, create it first
	dir, _ := filepath.Split(absPath)
	err = os.MkdirAll(dir, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
	if err != nil {
		panic(err)
	}
}

// GetFileContent returns required file content
func GetFileContent(projectid string, filePath string) []byte {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		panic(err)
	}
	absPath := filepath.Join(ROOT, projectName, filePath)

	// Read file content
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		panic(err)
	}
	return content
}

// GetFileStructure read file structure and return it
func GetFileStructure(projectid string) []entities.FileStructure {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		panic(err)
	}
	absPath := filepath.Join(ROOT, projectName)

	// Recurisively get file structure
	return dfs(absPath)

}

func dfs(path string) []entities.FileStructure {
	var structure []entities.FileStructure
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		tmp := entities.FileStructure{
			ID:         1,
			Name:       file.Name(),
			EditStatus: 0,
		}
		if file.IsDir() {
			tmp.Type = "dir"
			nextPath := filepath.Join(path, file.Name())
			tmp.Children = dfs(nextPath)
		} else {
			tmp.Type = "file"
		}
		structure = append(structure, tmp)
	}
	return structure
}
