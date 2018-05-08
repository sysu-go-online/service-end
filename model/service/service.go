package service

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sysu-go-online/service-end/model/entities"
)

// ROOT defines the root directory
var ROOT = "/home"

// GetProjectNameByID return project according to the given id
func GetProjectNameByID(id string) (string, error) {
	return "test", nil
}

// UpdateFileContent update content with given filepath and content
func UpdateFileContent(projectid string, username string, filePath string, content string, create bool, dir bool) error {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Update file, if the file not exists, judge accroding to the given param
	if create {
		if dir {
			err = os.Mkdir(filePath, os.ModeAppend)
		} else {
			err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
		}
	} else {
		err = ioutil.WriteFile(absPath, []byte(content), os.ModeAppend)
	}
	return err
}

// DeleteFile delete file accroding to the given path
func DeleteFile(projectid string, username string, filePath string) error {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Delete file
	err = os.RemoveAll(absPath)
	return err
}

// GetFileContent returns required file content
func GetFileContent(projectid string, username string, filePath string) ([]byte, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return nil, err
	}
	absPath := getFilePath(username, projectName, filePath)

	// Read file content
	content, err := ioutil.ReadFile(absPath)
	if err != nil {
		return nil, err
	}
	return content, nil
}

// GetFileStructure read file structure and return it
func GetFileStructure(projectid string, username string) ([]entities.FileStructure, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return nil, err
	}
	absPath := getFilePath(username, projectName, "/")

	// Recurisively get file structure
	return dfs(absPath)
}

var id = 1

func dfs(path string) ([]entities.FileStructure, error) {
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

func getFilePath(username string, projectName string, filePath string) string {
	return filepath.Join(ROOT, username, "go/src/github.com", username, projectName, filePath)
}
