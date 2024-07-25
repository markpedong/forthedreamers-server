package models

import (
	// "github.com/lib/pq"
	"gorm.io/plugin/soft_delete"
)

type Users struct {
	ID                   string                `json:"id" gorm:"primaryKey"`
	FirstName            string                `json:"first_name" validate:"required"`
	LastName             string                `json:"last_name" validate:"required"`
	Phone                string                `json:"phone" validate:"required"`
	Address              string                `json:"address" validate:"required"`
	Description          string                `json:"description" validate:"required"`
	Email                string                `json:"email" validate:"required"`
	Username             string                `json:"username"`
	Password             string                `json:"password"`
	ResumePDF            string                `json:"resume_pdf"`
	ResumeDocx           string                `json:"resume_docx"`
	IsResumeDownloadable int                   `json:"isResumeDownloadable" gorm:"default:0"`
	CreatedAt            int                   `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            int                   `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt            soft_delete.DeletedAt `json:"-"`
}
