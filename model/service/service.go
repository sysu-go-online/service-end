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
func UpdateFileContent(projectid string, filePath string, content string) error{
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return err
	}
	absPath := filepath.Join(ROOT, projectName, filePath)

	// Update file, if the file not exists, create it first
	dir, _ := filepath.Split(absPath)
	err = os.MkdirAll(dir, os.ModeAppend)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
	if err != nil {
		return err
	}
	return nil
}

// GetFileContent returns required file content
func GetFileContent(projectid string, filePath string) ([]byte, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return nil, err
	}
	absPath := filepath.Join(ROOT, projectName, filePath)

	// Read file content
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetFileStructure read file structure and return it
func GetFileStructure(projectid string) ([]entities.FileStructure, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return nil, err
	}
	absPath := filepath.Join(ROOT, projectName)

	// Recurisively get file structure
	return dfs(absPath)
}

var id = 1

func dfs(path string) ([]entities.FileStructure, error){
	var structure []entities.FileStructure
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		tmp := entities.FileStructure{
			ID:         id,
			Name:       file.Name(),
			EditStatus: 0,
		}
		id++
		if file.IsDir() {
			tmp.Type = "dir"
			nextPath := filepath.Join(path, file.Name())
			tmp.Children, err = dfs(nextPath)
			if err != nil {
				return nil, err
			}
		} else {
			tmp.Type = "file"
		}
		structure = append(structure, tmp)
	}
	return structure, nil
}
