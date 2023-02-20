package videos

import (
	"os"
	"sync"

	"github.com/HumXC/simple-douyin/model"
)

// video 和 cover 是完成后的视频和封面文件的路径，返回值表示是否删除文件
// 如果返回值为 true 将在函数调用后删除 video 和 cover 文件
type ButcherFinidhFunc = func(job Job, video, cover string, err error) (delete bool)

// Butcher 开启协程运行的函数
type ButcherWorkFunc = func(input string) (video, cover string, err error)

type VideoJobMan interface {
	// 添加一个任务，如果没有错误，传入的 job 将被补足 ID 等数据库字段
	Add(job model.VideoJob) uint
	// 获取所有未完成的 job
	Get() []model.VideoJob
	Rm(id uint) error
}

type Job struct {
	ID     uint
	Src    string
	Title  string
	UserID int64
	work   ButcherWorkFunc
	finish ButcherFinidhFunc
}

// 用于将压缩视频以及截取封面
// 工作过程是异步的，每完成一个任务都会调用对应的 WhenFinish
type Butcher struct {
	db      VideoJobMan
	tasks   []Job // 任务栈
	work    ButcherWorkFunc
	finish  ButcherFinidhFunc
	working chan struct{} // 正在进行的工作
	mu      sync.Mutex
}

// 创建一个 Butcher，maxJob 是同时工作的最大协程数量，在 work 完成后会调用 finish
func NewButcher(db VideoJobMan, maxJob int, finish ButcherFinidhFunc) *Butcher {
	return SNewButcher(db, maxJob, nil, finish)
}

func SNewButcher(db VideoJobMan, maxJob int, work ButcherWorkFunc, finish ButcherFinidhFunc) *Butcher {
	b := &Butcher{
		db:      db,
		work:    work,
		finish:  finish,
		working: make(chan struct{}, maxJob),
		tasks:   make([]Job, 0, 16),
	}
	if b.work == nil {
		b.work = defaultWork
	}
	// 从数据库获取未完成的工作
	w := b.db.Get()
	for _, v := range w {
		b.tasks = append(b.tasks, Job{
			ID:     v.ID,
			Src:    v.Src,
			UserID: v.UserID,
			Title:  v.Title,
		})
	}

	go b.do()
	return b
}

// 添加一个任务，input 是需要处理的文件
func (b *Butcher) Add(input, title string, userID int64) {
	b.SAdd(input, title, userID, nil, nil)
}

// 添加一个任务，input 是需要处理的文件，finish 是完成后调用的函数
func (b *Butcher) SAdd(input, title string, userID int64, work ButcherWorkFunc, finish ButcherFinidhFunc) {
	job := model.VideoJob{
		Src:    input,
		UserID: userID,
		Title:  title,
	}
	id := b.db.Add(job)

	b.tasks = append(b.tasks, Job{
		ID:     id,
		Src:    input,
		Title:  title,
		UserID: userID,
		work:   work,
		finish: finish,
	})
	go b.do()
}

// 开始对视频进行压缩工作
// 启用多个协程同时进行，协程的数量取决于 cap(working)
// 如果 work 运行没有错误，则从数据库中移除任务表示已经完成
func (b *Butcher) do() {
	if !b.mu.TryLock() {
		return
	}
	if len(b.tasks) == 0 {
		b.mu.Unlock()
		return
	}
	for len(b.tasks) != 0 {
		var jobs []Job
		if len(b.tasks) <= cap(b.working) {
			jobs = b.tasks[:]
			b.tasks = b.tasks[1:]
		} else {
			jobs = b.tasks[:cap(b.working)]
			b.tasks = b.tasks[cap(b.working):]
		}
		for _, job := range jobs {
			b.working <- struct{}{}
			go func(job Job) {
				if job.work == nil {
					job.work = b.work
				}
				if job.finish == nil {
					job.finish = b.finish
				}
				v, c, err := job.work(job.Src)
				if job.finish != nil && job.finish(job, v, c, err) {
					_ = os.Remove(job.Src)
					_ = os.Remove(v)
					_ = os.Remove(c)
				}
				// 如果出错就不会移除任务
				if err == nil {
					b.db.Rm(job.ID)
				}
				b.db.Rm(job.ID)
				_ = <-b.working
			}(job)
		}
	}
	b.mu.Unlock()
}

func defaultWork(input string) (string, string, error) {
	video, err := SmallVideoWithFfmpeg(input)
	if err != nil {
		return "", "", err
	}
	cover, err := CutVideoWithFfmpeg(video)
	if err != nil {
		return "", "", err
	}
	return video, cover, nil
}
