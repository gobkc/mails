package mails

import (
	"reflect"
	"sync"
)

type MailType byte

const (
	MailTypeHtml MailType = iota
	MailTypeText
)

var types = [][]byte{
	{67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 58, 32, 116, 101, 120, 116, 47, 104, 116, 109, 108, 59, 32, 99, 104, 97, 114, 115, 101, 116, 61, 85, 84, 70, 45, 56},
	{67, 111, 110, 116, 101, 110, 116, 45, 84, 121, 112, 101, 58, 32, 116, 101, 120, 116, 47, 112, 108, 97, 105, 110, 59, 32, 99, 104, 97, 114, 115, 101, 116, 61, 85, 84, 70, 45, 56},
}

func (t MailType) Bytes() []byte {
	//if int(t) >= len(types) {
	if int(t) >= 2 {
		return []byte{}
	}
	return types[t]
}

type Email interface {
	SetEnv(userName, password, mailServer string, mailType ...MailType) error
	SendToMail(toEmail, subject, body string) error
	SendOTPToMail(toEmail, subject, body string) error
	VerifyOTP(email, code string) error
}

var cache = make(map[any]any)
var l sync.Mutex

func EmailFactory[T any]() (t *T) {
	target := reflect.TypeOf(t)
	l.Lock()
	defer l.Unlock()
	find, ok := cache[target]
	if ok {
		return find.(*T)
	}
	v := new(T)
	cache[target] = v
	return v
}
