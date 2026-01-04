package models

import (
	"golang-http-patch/patch"

	"gorm.io/gorm"
)

// User model
type User struct {
	ID     uint    `json:"id" gorm:"primaryKey"`
	Name   string  `json:"name" gorm:"not null"`
	Email  string  `json:"email" gorm:"uniqueIndex;not null"`
	Age    int     `json:"age" gorm:"not null"`
	Phone  *string `json:"phone" gorm:"type:varchar(20)"` // nullable - can be null
	Active bool    `json:"active" gorm:"default:true"`
	Bio    string  `json:"bio" gorm:"type:text"`                        // optional text field
	Role   string  `json:"role" gorm:"type:varchar(20);default:'user'"` // enum-like: admin, user, guest
	Score  float64 `json:"score" gorm:"type:decimal(10,2);default:0"`   // numeric field
}

// DTOs for request validation
type CreateUserDTO struct {
	Name   string  `json:"name" validate:"required,min=2,max=100"`
	Email  string  `json:"email" validate:"required,email"`
	Age    int     `json:"age" validate:"required,gte=0,lte=150"`
	Phone  *string `json:"phone" validate:"omitempty,min=10,max=20"` // optional, nullable
	Active bool    `json:"active"`                                   // optional, defaults to true
	Bio    string  `json:"bio" validate:"omitempty,max=500"`         // optional text
	Role   string  `json:"role" validate:"omitempty,oneof=admin user guest"`
	Score  float64 `json:"score" validate:"omitempty,gte=0,lte=100"`
}

type UpdateUserDTO struct {
	Name   string  `json:"name" validate:"required,min=2,max=100"`
	Age    int     `json:"age" validate:"required,gte=0,lte=150"`
	Phone  *string `json:"phone" validate:"omitempty,min=10,max=20"`
	Active bool    `json:"active"`
	Bio    string  `json:"bio" validate:"omitempty,max=500"`
	Role   string  `json:"role" validate:"required,oneof=admin user guest"`
	Score  float64 `json:"score" validate:"required,gte=0,lte=100"`
	// Email is immutable and cannot be updated after creation
}

type PatchUserDTO struct {
	// String field: can be unset (ignore), null (remove), or value (update)
	Name patch.Optional[string] `json:"name" validate:"opt=min=2;max=100"`

	// Numeric field: validate range only when present
	Age patch.Optional[int] `json:"age" validate:"opt=gte=0;lte=150"`

	// Nullable string field: can be unset (ignore), null (remove), or value (update)
	Phone patch.Optional[string] `json:"phone" validate:"opt=min=10;max=20"`

	// Boolean field: can be unset (ignore), null (remove), or value (update)
	Active patch.Optional[bool] `json:"active"`

	// Optional text field: can be unset (ignore), null (remove), or value (update)
	Bio patch.Optional[string] `json:"bio" validate:"opt=max=500"`

	// Enum-like field: can be unset (ignore), null (remove), or value (update)
	Role patch.Optional[string] `json:"role" validate:"opt=oneof=admin user guest"`

	// Float field: can be unset (ignore), null (remove), or value (update)
	Score patch.Optional[float64] `json:"score" validate:"opt=gte=0;lte=100"`
	// Email is immutable and cannot be updated after creation
}

// AutoMigrate runs database migrations for User model
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
