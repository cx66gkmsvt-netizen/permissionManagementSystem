package task

import (
	"log"
	"user-center/internal/service"

	"github.com/robfig/cron/v3"
)

var c *cron.Cron

// InitCronJobs initializes and starts the scheduled tasks
func InitCronJobs() {
	// Create a new cron job runner
	c = cron.New()

	attendanceSvc := service.NewAttendanceService()

	// 每天 00:00:00 触发在班状态重置为休班
	_, err := c.AddFunc("0 0 * * *", func() {
		log.Println("Running scheduled task: Resetting daily attendance to off-duty")
		if err := attendanceSvc.ResetDailyAttendance(); err != nil {
			log.Printf("Error resetting daily attendance via cron: %v", err)
		}
	})

	if err != nil {
		log.Fatalf("Failed to add cron job for resetting daily attendance: %v", err)
	}

	c.Start()
	log.Println("Cron jobs initialized and started")
}
