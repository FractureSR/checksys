package main

import (
	"net/http"

	"checksys2.0/apratus"
	"checksys2.0/handler"
	"checksys2.0/threads"
)

func main() {
	apratus.PushInfoDeliver = make(chan apratus.PushInfo)           //创建传输教师发布课程信息的通道的实例
	apratus.RestudentInfoDeliver = make(chan apratus.RestudentInfo) //创建传输学生端返回信息的通道的实例
	apratus.LectureInfo = make(map[string]apratus.PushInfo)         //创建存储课程信息的映射
	apratus.Recorders = make(map[string]apratus.Recorder)           //创建存储记录器的映射
	go threads.Core(apratus.PushInfoDeliver)                        //打开核心协程，该协程用于集中接收教师发布的课程信息并处理
	go threads.ViceCore(apratus.RestudentInfoDeliver)               //打开复核心协程，该协程主要处理学生返回信息

	mux := http.NewServeMux() //多路复用

	//绑定访问标识符和处理器函数
	mux.HandleFunc("/test", handler.TestHandler)
	mux.HandleFunc("/studentAccess", handler.StudentAccess)
	mux.HandleFunc("/teacherAccess", handler.TeacherAccess)
	mux.HandleFunc("/studentlogin", handler.StudentLogin)
	mux.HandleFunc("/studentInfo", handler.StudnetInfo)
	mux.HandleFunc("/studentCheck", handler.StudentCheck)
	mux.HandleFunc("/teacherlogin", handler.TeacherLogin)
	mux.HandleFunc("/check", handler.Check)
	mux.HandleFunc("/lectureInfo", handler.LectureInfo)
	//静态文件系统用于提供静态文件支援
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files)) //对一些url进行处理

	http.ListenAndServe("localhost:80", mux) //监听
}
