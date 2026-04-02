package scheduler

import (
	"log"
	"time"
	"workflow-system/internal/repository"
	"workflow-system/internal/service/engine"
)

// TimeoutScheduler 超时调度器
type TimeoutScheduler struct {
	instanceRepo *repository.InstanceRepository
	engine       *engine.WorkflowEngine
	interval     time.Duration // 检查间隔
	stopCh       chan struct{}
}

// NewTimeoutScheduler 创建超时调度器
func NewTimeoutScheduler(instanceRepo *repository.InstanceRepository, wfEngine *engine.WorkflowEngine, interval time.Duration) *TimeoutScheduler {
	return &TimeoutScheduler{
		instanceRepo: instanceRepo,
		engine:       wfEngine,
		interval:     interval,
		stopCh:       make(chan struct{}),
	}
}

// Start 启动调度器
func (s *TimeoutScheduler) Start() {
	log.Printf("Timeout scheduler started, interval: %v", s.interval)
	ticker := time.NewTicker(s.interval)
	go func() {
		// 立即执行一次
		s.checkOverdueTasks()
		for {
			select {
			case <-ticker.C:
				s.checkOverdueTasks()
			case <-s.stopCh:
				ticker.Stop()
				log.Println("Timeout scheduler stopped")
				return
			}
		}
	}()
}

// Stop 停止调度器
func (s *TimeoutScheduler) Stop() {
	close(s.stopCh)
}

// checkOverdueTasks 检查并处理超时任务
func (s *TimeoutScheduler) checkOverdueTasks() {
	now := time.Now().Unix()
	tasks, err := s.instanceRepo.GetOverduePendingTasks(now)
	if err != nil {
		log.Printf("Failed to get overdue tasks: %v", err)
		return
	}

	for _, task := range tasks {
		log.Printf("Processing overdue task: %d, deadline was: %d", task.ID, *task.DeadlineAt)
		// 超时自动通过（或者可以根据配置设置为驳回）
		if err := s.engine.ProcessApproval(task.ID, 0, "approve", "超时自动通过"); err != nil {
			log.Printf("Failed to process overdue task %d: %v", task.ID, err)
		} else {
			log.Printf("Overdue task %d auto-approved", task.ID)
		}
	}
}
