# About Mails
A library that simplifies sending mail using "net/smtp"

### Contributing
You can commit PR to this repository

### How to get it?
````
go get -u github.com/gobkc/mails
````

### Quick start
````
package main

import (
	"github.com/gobkc/mails"
	"log"
)

func main() {
	mail := mails.EmailFactory[mails.DefaultEmail]()
	if err := mail.SetEnv("example@163.com", "examplepassword", "smtp.163.com:25"); err != nil {
		log.Println(err.Error())
		return
	}
	if err := mail.SendOTPToMail("example_send_to_mail@qq.com", "test otp", ""); err != nil {
		log.Println(err.Error())
		return
	}
	if err := mail.SendToMail("example_send_to_mail@qq.com", "test title", "test content"); err != nil {
		log.Println(err.Error())
		return
	}
}
````

Note:The sender's SMTP function must be enabled

### License
© Gobkc, 2022~time.Now

Released under the Apache License