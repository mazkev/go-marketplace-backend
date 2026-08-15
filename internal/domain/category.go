package domain

import "context"

type Category struct {
	ID        uint       `gorm:"primaryKey" json:"id"`
	ParentID  *uint      `gorm:"index" json:"parent_id"`
	Name      string     `gorm:"size:255;not null" json:"name"`
	Slug      string     `gorm:"size:255;uniqueIndex;not null" json:"slug"`
	Children  []Category `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetAll(ctx context.Context) ([]Category, error)
	GetByID(ctx context.Context, id uint) (*Category, error)
}
