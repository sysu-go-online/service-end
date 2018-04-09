package controller

import (
	"github.com/sysu-go-online/service-end/model/entities"
)

type Project struct {
	UserID      int    `json:"userid"`
	ProjectName string `json:"projectName"`
	Type        int    `json:"type"`
}

func NewProject(raw Project) entities.Project {
	ret := entities.Project{
		ProjectName: raw.ProjectName,
		UserID:      raw.UserID,
		Language:    raw.Type,
	}
	return ret
}
