package model

type BlacklistMember struct {
	Id   int64  `gorm:"type:INT UNSIGNED NOT NULL DEFAULT 0;index;comment:identifier"`
	Type string `gorm:"type:VARCHAR(20) NOT NULL;comment:private or group"`
}
