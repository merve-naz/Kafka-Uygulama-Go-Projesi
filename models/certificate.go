package models

import "time"

type Certificate struct {
	RegistrationUUID string    `json:"registration_uid" validate:"required" binding:"omitempty"`
	CompletedAt      time.Time `json:"completed_at" validate:"required" binding:"omitempty"`
	RegisteredAt     time.Time `json:"registered_at" validate:"required" binding:"omitempty"`
	Name             string    `json:"name" validate:"required" binding:"omitempty"`
	Slug             string    `json:"slug" validate:"required" binding:"omitempty"`
	Avatar           string    `json:"avatar" validate:"required" binding:"omitempty"`
	Title            string    `json:"title" validate:"required" binding:"omitempty"`
	UrlSlug          string    `json:"url_slug" validate:"required" binding:"omitempty"`
	BadgeOwner       string    `json:"badge_owner" validate:"required" binding:"omitempty"`
	BadgeTitle       string    `json:"badge_title" validate:"required" binding:"omitempty"`
}
