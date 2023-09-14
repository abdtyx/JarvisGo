package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedAt time.Time      `gorm:"type:DATETIME(3) NOT NULL;comment:创建时间"`
	UpdatedAt time.Time      `gorm:"type:DATETIME(3) NOT NULL;comment:更新时间"`
	DeletedAt gorm.DeletedAt `gorm:"type:DATETIME(3) NULL;index;comment:删除时间"`
}

type Blacklist struct {
	PK      int64  `gorm:"type:BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT;primaryKey;comment:主键"`
	Id      uint   `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
	Type    string `gorm:"type:VARCHAR(20) NOT NULL;comment:user or group"`
	Comment string `gorm:"type:VARCHAR(20) NOT NULL;comment:comment"`

	BaseModel
}

func (b *Blacklist) String() string {
	return fmt.Sprintf("Id: %v, Type: %v, Comment: %v\n", b.Id, b.Type, b.Comment)
}

func (Blacklist) TableName() string {
	return "blacklists"
}

type Jeminder struct {
	Id uint `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;primaryKey;comment:identifier"`

	BaseModel
}

func (Jeminder) TableName() string {
	return "jeminders"
}
