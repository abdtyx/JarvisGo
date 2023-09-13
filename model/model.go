package model

import "fmt"

type Blacklist struct {
	Id      uint   `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
	Type    string `gorm:"type:VARCHAR(20) NOT NULL;comment:user or group"`
	Comment string `gorm:"type:VARCHAR(20) NOT NULL;comment:comment"`
}

func (b *Blacklist) String() string {
	return fmt.Sprintf("Id: %v, Type: %v, Comment: %v\n", b.Id, b.Type, b.Comment)
}

type Jeminder struct {
	Id uint `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
}
