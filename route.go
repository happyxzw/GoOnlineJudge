package main

import (
	"GoOnlineJudge/controller"
	"GoOnlineJudge/controller/admin"
	"GoOnlineJudge/controller/contest"
	"net/http"
	"reflect"
	"strings"
)

// Page

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		c := &controller.NewsController{}
		m := "List"
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

func newsHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := &controller.NewsController{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

func problemHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := &controller.ProblemController{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := &controller.StatusController{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

func ranklistHandler(w http.ResponseWriter, r *http.Request) {
	c := &controller.RanklistController{}
	m := "Index"
	rv := getReflectValue(w, r)
	callMethod(c, m, rv)
}

func contestlistHandler(w http.ResponseWriter, r *http.Request) {
	c := &controller.ContestController{}
	m := "List"
	rv := getReflectValue(w, r)
	callMethod(c, m, rv)
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := &controller.UserController{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

func helpHandler(w http.ResponseWriter, r *http.Request) {
	p := strings.Trim(r.URL.Path, "/")
	s := strings.Split(p, "/")
	if l := len(s); l >= 2 {
		c := %controller.HelpController{}
		m := strings.Title(s[1])
		rv := getReflectValue(w, r)
		callMethod(c, m, rv)
	}
}

// Contest
func contestHandler(w http.ResponseWriter, r *http.Request) {
	c := &contest.ContestUserContorller{}
	m := "Register"
	rv := getReflectValue(w, r)
	callMethod(c, m, rv)
}

// Admin

func adminHandler(w http.ResponseWriter, r *http.Request) {
	c := &admin.AdminUserController{}
	m := "Register"
	rv := getReflectValue(w, r)
	callMethod(c, m, rv)
}

// Common

func callMethod(c interface{}, m string, rv []reflect.Value) {
	rc := reflect.ValueOf(c)
	rm := rc.MethodByName(m)
	rm.Call(rv)
}

func getReflectValue(w http.ResponseWriter, r *http.Request) (rv []reflect.Value) {
	rw := reflect.ValueOf(w)
	rr := reflect.ValueOf(r)
	rv = []reflect.Value{rw, rr}
	return
}
