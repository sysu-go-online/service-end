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
func GetFileStructure(projectid string, username string) (*entities.FileStructure, error) {
	// Get absolute path
	projectName, err := GetProjectNameByID(projectid)
	if err != nil {
		return nil, err
	}
	absPath := getFilePath(username, projectName, "/")

	// Recurisively get file structure
	s, err := dfs(absPath, 0)
	if err != nil {
		return nil, err
	}
	// Add root content
	root := entities.FileStructure{
		ID: 1,
		Name: projectName,
		Type: "dir",
		Children: s,
		Root: true,
		IsSelected: true,
	}
	return &root, nil
}

var id = 1

func dfs(path string, depth int) ([]entities.FileStructure, error) {
	var structure []entities.FileStructure
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		tmp := entities.FileStructure{
			ID:         id,
			Name:       file.Name(),
			Root:       false,
			IsSelected: false,
		}
		id++
		if file.IsDir() {
			tmp.Type = "dir"
			nextPath := filepath.Join(path, file.Name())
			tmp.Children, err = dfs(nextPath, depth+1)
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
	return filepath.Join(ROOT, username, "go/src/github.com", projectName, filePath)
}

// TODO: Get user information
func GetUserInformation(username string) entities.UserInfo {
	return entities.UserInfo{}
}

// AddUser Insert user information
func AddUser(e entities.UserInfo) error {
	_, err := entities.Engine.InsertOne(e)
	return err
}

// UpdateAccessToken update token in the db
func UpdateAccessToken(token string) error {
	e := entities.UserInfo{
		Token: token,
	}
	_, err := entities.Engine.Update(e)
	return err
}
