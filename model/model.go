package model

type Blacklist struct {
	Id   uint   `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
	Type string `gorm:"type:VARCHAR(20) NOT NULL;comment:private or group"`
}

type Jeminder struct {
	Id uint `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
}
