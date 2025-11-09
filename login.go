package main

import (
	"net/http"

	"github.com/west2-online/jwch"
	"golang.org/x/crypto/bcrypt"
)

type session struct {
	Identity string
	Cookies  []*http.Cookie
	Password string
}

var sessionCache = make(map[string]session)

func Login(uid string, password string) (*jwch.Student, error) {
	s, exists := sessionCache[uid]
	var stu *jwch.Student
	if exists && checkPasswordHash(password, s.Password) {
		stu = jwch.NewStudent().WithLoginData(s.Identity, s.Cookies)
		err := stu.Login()
		if err == nil {
			return stu, nil
		}
	}
	stu = jwch.NewStudent().WithUser(uid, password)
	err := stu.Login()
	if err != nil {
		return nil, err
	}

	identity, cookies, err := stu.GetIdentifierAndCookies()

	if err != nil {
		return nil, err
	}

	sessionCache[uid] = session{
		Identity: identity,
		Cookies:  cookies,
		Password: hashPassword(password),
	}

	return stu, nil
}

func hashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(hash))
	return err == nil
}
