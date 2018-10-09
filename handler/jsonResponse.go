package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"checksys2.0/apratus"
	_ "github.com/go-sql-driver/mysql"
)

func StudnetInfo(w http.ResponseWriter, r *http.Request) {
	//打开数据库
	db, err := sql.Open("mysql", "root:Neu1923#@tcp(localhost:3306)/checksys_data?charset=utf8")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var stuName string
	var stuCollege string
	var stuClass string
	var stuMajor string
	//通过id取得学生的信息
	err = db.QueryRow("SELECT name,college,major,class FROM studentinfo where stuid=?", r.PostFormValue("id")).Scan(&stuName, &stuCollege, &stuMajor, &stuClass)
	if err != nil {
		panic(err)
	}

	//构建学生信息实例
	AstudentInfo := apratus.StudentInfo{
		Name:    stuName,
		College: stuCollege,
		Major:   stuMajor,
		Class:   stuClass,
		Id:      r.PostFormValue("id"),
	}
	//编码
	output, err := json.MarshalIndent(&AstudentInfo, "", "\t\t")
	if err != nil {
		fmt.Println("Make stu Json error.")
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	//发送
	w.Write(output)
}

func LectureInfo(w http.ResponseWriter, r *http.Request) {
	var lectureList []apratus.PushInfo
	//取得课程数据
	mutex.Lock()
	{
		for _, value := range apratus.LectureInfo {
			lectureList = append(lectureList, value)
		}
	}
	mutex.Unlock()
	//编码并发送
	output, err := json.MarshalIndent(&lectureList, "", "\t\t")
	if err != nil {
		fmt.Println("Make lecture Json error.")
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
}
