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
	/*
		TODO:
		- add notification low stock to cron job hourly
	*/
	if err != nil {
		log.Fatal("can't init service context %w", err)
	}

	_, err = c.AddFunc("@weekly", func() {
		ctx.Ctl.UserHandler.UpdateShiftSchedule()
	})

	if err != nil {
		fmt.Println("Error update shift schedule from cron:", err)
		return
	}
	// Menjalankan cron
	c.Start()
	// Menunggu agar main tidak langsung selesai
	fmt.Println("Cron job berjalan. Tekan CTRL+C untuk keluar.")
	select {}
}
