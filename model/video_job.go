package model

import (
	"sync"

	"gorm.io/gorm"
)

type VideoJob struct {
	gorm.Model
	Src    string // 待处理视频的源文件
	Title  string
	UserID int64
}

type VideoJobMan struct {
	db *gorm.DB
	mu sync.Mutex
}

// 添加一个任务，如果没有错误，传入的 job 将被补足 ID 等数据库字段
func (v *VideoJobMan) Add(job VideoJob) uint {
	v.mu.Lock()
	defer v.mu.Unlock()
	_ = v.db.Create(&job).Error
	_ = v.db.Select("id").Last(&job).Error
	return job.ID
}

// 从数据库获取所有未完成的 job
func (v *VideoJobMan) Get() []VideoJob {
	jobs := make([]VideoJob, 0, 16)
	_ = v.db.Find(&jobs).Error
	return jobs
}

func (v *VideoJobMan) Rm(id uint) error {
	return v.db.Model(VideoJob{}).Delete("id=?", id).Error
}
