package sqldb

import (
	"sync"

	"github.com/HumXC/simple-douyin/model"
	"gorm.io/gorm"
)

type videoJobMan struct {
	db *gorm.DB
	mu sync.Mutex
}

// 添加一个任务，如果没有错误，传入的 job 将被补足 ID 等数据库字段
func (v *videoJobMan) Add(job model.VideoJob) uint {
	v.mu.Lock()
	defer v.mu.Unlock()
	_ = v.db.Create(&job).Error
	_ = v.db.Select("id").Last(&job).Error
	return job.ID
}

// 从数据库获取所有未完成的 job
func (v *videoJobMan) Get() []model.VideoJob {
	jobs := make([]model.VideoJob, 0, 16)
	_ = v.db.Find(&jobs).Error
	return jobs
}

func (v *videoJobMan) Rm(id uint) error {
	return v.db.Model(model.VideoJob{}).Delete("id=?", id).Error
}
