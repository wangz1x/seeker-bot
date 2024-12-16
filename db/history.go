// history - 2024/12/16
// Author: wangzx
// Description:

package db

import "time"

type History struct {
	ID        int       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserID    string    `gorm:"column:user_id;type:varchar(500);not null" json:"user_id"`
	Role      string    `gorm:"column:role;type:varchar(100);not null" json:"role"`
	Content   string    `gorm:"column:content;type:text" json:"content"`
	CreatedAt time.Time `gorm:"column:created_at;type:timestamp;default:CURRENT_TIMESTAMP" json:"created_at"`
	MsgID     string    `gorm:"column:msg_id;type:varchar(500);not null" json:"msg_id"`
}

// TableName 指定表名
func (History) TableName() string {
	return "history"
}
