package main

import (
	mail "github.com/xhit/go-simple-mail/v2"
	"github.com/yusufelyldrm/reservation/internal/models"
	_ "net/mail"
	"time"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMsg(msg)
		}
	}()
}

func sendMsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		app.ErrorLog.Println(err)
		return
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, string(m.Content))
	err = email.Send(client)
	if err != nil {
		app.ErrorLog.Println(err)
		return
	} else {
		app.InfoLog.Println("Email sent")
	}
}
