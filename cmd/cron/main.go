package main

import (
	"fmt"
	"log"
	"project/infra"

	"github.com/robfig/cron/v3"
)

func main() {
	// Membuat instance cron
	c := cron.New()
	ctx, err := infra.NewServiceContext()
	if err != nil {
		log.Fatal("can't init service context %w", err)
	}

	if err = cronJobs(c, ctx); err != nil {
		return
	}

	// Menjalankan cron
	c.Start()
	// Menunggu agar main tidak langsung selesai
	fmt.Println("Cron job berjalan. Tekan CTRL+C untuk keluar.")
	select {}
}

func cronJobs(c *cron.Cron, ctx *infra.ServiceContext) error {
	if err := updateShiftSchedule(c, ctx); err != nil {
		return err
	}

	if err := sendNotificationLowStock(c, ctx); err != nil {
		return err
	}

	if err := addBestSeller(c, ctx); err != nil {
		return err
	}

	return nil
}

func updateShiftSchedule(c *cron.Cron, ctx *infra.ServiceContext) error {
	if _, err := c.AddFunc("@weekly", func() {
		ctx.Ctl.UserHandler.UpdateShiftSchedule()
	}); err != nil {
		fmt.Println("Error update shift schedule from cron:", err)
		return err
	}

	return nil
}

func sendNotificationLowStock(c *cron.Cron, ctx *infra.ServiceContext) error {
	if _, err := c.AddFunc("* * * * *", func() {
		ctx.Ctl.NotificationHandler.SendNotificationLowStock()
	}); err != nil {
		fmt.Println("Error sending notification low stock from cron:", err)
		return err
	}

	return nil
}

func addBestSeller(c *cron.Cron, ctx *infra.ServiceContext) error {
	if _, err := c.AddFunc("0 1 * * *", func() {
		ctx.Ctl.RevenueHandler.AddDailyBestSeller(ctx.Cfg.ProfitMargin)
	}); err != nil {
		fmt.Println("Error adding best seller from cron:", err)
		return err
	}

	return nil
}
