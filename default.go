package mails

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"html/template"
	"net/smtp"
	"strings"
	"time"
)

type DefaultEmail struct {
	user     string
	pass     string
	server   string
	fullServ string
	ct       []byte
}

func (d *DefaultEmail) SetEnv(userName, password, mailServer string, mailType ...MailType) error {
	d.fullServ = mailServer
	serverList := strings.Split(mailServer, ":")
	if len(serverList) < 2 {
		return fmt.Errorf(`SetEnv:server address error,like "smtp.163.com:25"`)
	}
	d.server = serverList[0]
	if len(mailType) >= 1 {
		d.ct = mailType[0].Bytes()
	} else {
		d.ct = MailTypeHtml.Bytes()
	}
	d.user = userName
	d.pass = password
	return nil
}

func (d *DefaultEmail) SendToMail(to, subject, body string) error {
	var msg bytes.Buffer
	msg.WriteString("To: ")
	msg.WriteString(to)
	msg.WriteString("\r\nFrom: ")
	msg.WriteString(d.user)
	msg.WriteString("\r\nSubject: ")
	msg.WriteString(subject)
	msg.WriteString("\r\n")
	msg.Write(d.ct)
	msg.WriteString("\r\n\r\n")
	msg.WriteString(body)
	auth := smtp.PlainAuth("", d.user, d.pass, d.server)
	sendTo := strings.Split(to, ";")
	err := smtp.SendMail(d.fullServ, auth, d.user, sendTo, msg.Bytes())
	return err
}

type OTPParam struct {
	Code     string
	UserName string
}

type otpCache struct {
	Code       string
	UserName   string
	Email      string
	HourTimes  uint8
	UpdateTime time.Time
}

func (d *DefaultEmail) SendOTPToMail(to, subject, body string) error {
	l.Lock()
	defer l.Unlock()
	otp := strings.ToUpper(Random(4))
	if body == `` {
		body = otpTemplate
	}
	param := OTPParam{
		Code: otp,
	}
	if find, ok := cache[to]; ok {
		c := find.(otpCache)
		c.Code = otp
		c.Email = to
		c.UpdateTime = time.Now()
		if c.UpdateTime.Add(1*time.Hour).Unix() < time.Now().Unix() {
			c.HourTimes = 1
		} else {
			if c.HourTimes >= 3 {
				return fmt.Errorf(`SendOTPToMail:HourTimes must lte 3`)
			}
			c.HourTimes++
		}
		cache[to] = c
	} else {
		cache[to] = otpCache{
			Code:       otp,
			Email:      to,
			HourTimes:  1,
			UpdateTime: time.Now(),
		}
	}
	c, err := template.New("member").Parse(body)
	if err != nil {
		return fmt.Errorf("SendOTPToMail:Parse:%w", err)
	}
	var buf bytes.Buffer
	err = c.Execute(&buf, param)
	if err != nil {
		return fmt.Errorf("SendOTPToMail:Execute:%w", err)
	}
	result := buf.String()
	return d.SendToMail(to, subject, result)
}

func (d *DefaultEmail) VerifyOTP(email, code string) error {
	l.Lock()
	defer l.Unlock()
	find, ok := cache[email]
	if !ok {
		return fmt.Errorf("VerifyOTP:EmailNotExists")
	}
	otp := find.(otpCache)
	if otp.Code != code {
		return fmt.Errorf("VerifyOTP:OTPCodeFailed")
	}
	if otp.UpdateTime.Add(1*time.Hour).Unix() < time.Now().Unix() {
		return fmt.Errorf("VerifyOTP:TimedOut")
	}
	return nil
}

func Random(strLen int) (str string) {
	var (
		randByte  = make([]byte, strLen)
		formatStr []string
		out       []interface{}
		byteHalf  uint8 = 127
	)
	if strLen == 0 {
		return
	}
	rand.Read(randByte)
	for _, b := range randByte {
		if b > byteHalf {
			formatStr = append(formatStr, "%X")
		} else {
			formatStr = append(formatStr, "%x")
		}
		out = append(out, b)
	}
	if str = fmt.Sprintf(strings.Join(formatStr, ""), out...); len(str) > strLen {
		str = str[:strLen]
	}
	return
}

var otpTemplate = `
<!DOCTYPE html>

<html lang="en">
<head>
    <meta name="viewport" content="width=device-width" />
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
    <title>OTP</title>
</head>
<body>
<section>
    <p>We have sent you 4 digits OTP.</p>
    <p>Verification Code: {{ .Code }}</p>
    <div class="divide-line"></div>
    <p>
        It will expire in 5 minutes, please use it as soon as possible
    </p>
</section>
</body>
</html>
`
