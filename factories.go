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

func Factory(emailImp Email) Email {
	target := reflect.TypeOf(emailImp)
	l.Lock()
	defer l.Unlock()
	find, ok := cache[target]
	findValue := reflect.ValueOf(find)
	if ok {
		typeOf := reflect.TypeOf(emailImp)
		valueOf := reflect.ValueOf(emailImp)
		if typeOf.Kind() == reflect.Pointer {
			typeOf = typeOf.Elem()
			valueOf = valueOf.Elem()
			findValue = findValue.Elem()
		}
		valueOf.Set(findValue)
		return emailImp
	}
	cache[target] = emailImp
	return emailImp
}
