package handler

import (
	"html/template"
	"net/http"
)

//学生登陆界面
func StudentAccess(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/StudentAccess.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}

//教师登陆界面
func TeacherAccess(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("templates/TeacherAccess.html")
	if err != nil {
		panic(err)
	}
	t.Execute(w, nil)
}
