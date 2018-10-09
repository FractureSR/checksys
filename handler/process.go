package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"checksys2.0/apratus"

	_ "github.com/go-sql-driver/mysql"
)

func StudentLogin(w http.ResponseWriter, r *http.Request) {
	//此处从浏览器接受了一个表单，内容是用户名和密码
	//打开数据库根据学号搜索密码
	db, err := sql.Open("mysql", "root:Neu1923#@tcp(localhost:3306)/checksys_data?charset=utf8")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var password = "123456789!@#$%^&*("
	err = db.QueryRow("SELECT password FROM studentinfo where stuid=?", r.PostFormValue("id")).Scan(&password)
	if err != nil {
		panic(err)
	}
	//匹配失败
	if password != r.PostFormValue("password") {
		t, terr := template.ParseFiles("templates/StudentAccess.html")
		if terr != nil {
			panic(terr)
		}
		t.Execute(w, "User do not exist or wrong password!")
		//匹配成功
	} else {
		t, terr := template.ParseFiles("templates/StudentSignIn.html")
		if terr != nil {
			panic(terr)
		}
		if err != nil {
			panic(err)
		}
		//设置cookie，包含了先前填入的学号
		cookieId := http.Cookie{
			Name:     "stuId",
			Value:    r.PostFormValue("id"),
			HttpOnly: false,
		}
		http.SetCookie(w, &cookieId)

		t.Execute(w, nil)
	}
}

func TeacherLogin(w http.ResponseWriter, r *http.Request) { //r中接收了一个表单，包含id和密码
	db, err := sql.Open("mysql", "root:Neu1923#@tcp(localhost:3306)/checksys_data?charset=utf8") //打开数据库
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var password = "123456789!@#$%^&*("
	var teacherName string
	var teacherCollege string
	err = db.QueryRow("SELECT password,nickname,college FROM teacherinfo WHERE teacherId=?", r.PostFormValue("teacherId")).Scan(&password, &teacherName, &teacherCollege)
	fmt.Println(r.PostFormValue("teacherId") + "\n" + r.PostFormValue("password") + "\n" + password)
	if password != r.PostFormValue("password") { // 匹配失败时重新显示TeacherAcess界面，并加入提示错误的信息
		t, terr := template.ParseFiles("templates/TeacherAccess.html")
		if terr != nil {
			panic(terr)
		}
		t.Execute(w, "User do not exist or wrong password!")
	} else { //匹配成功
		t, terr := template.ParseFiles("templates/TeacherPush.html")
		if terr != nil {
			panic(terr)
		}
		rows, err := db.Query("SELECT sunday,monday,tuesday,wednesday,thursday,friday,saturday FROM timetable WHERE belonging=?", r.PostFormValue("teacherId")) //取出与id匹配的数据库中的教师课表的rows
		if err != nil {
			panic(err)
		}
		var Atimetable apratus.Timetable
		for rows.Next() { //当下一条记录不为空时一直执行
			var sunday, monday, tuesday, wednesday, thursday, friday, saturday string
			err := rows.Scan(&sunday, &monday, &tuesday, &wednesday, &thursday, &friday, &saturday) //将取出的信息填入各个weekday
			if err != nil {
				panic(err)
			}
			//将读到的信息填入时间表
			Atimetable.Sunday = append(Atimetable.Sunday, sunday)
			Atimetable.Monday = append(Atimetable.Monday, monday)
			Atimetable.Tuesday = append(Atimetable.Tuesday, tuesday)
			Atimetable.Wednesday = append(Atimetable.Wednesday, wednesday)
			Atimetable.Thursday = append(Atimetable.Thursday, thursday)
			Atimetable.Friday = append(Atimetable.Friday, friday)
			Atimetable.Saturday = append(Atimetable.Saturday, saturday)
		}

		//对timetable编码
		output, err := json.MarshalIndent(&Atimetable, "", "\t\t")
		if err != nil {
			panic(err)
		}

		//写入文件，文件名对应于教师的id唯一
		var writePath = "./public/temp/" + r.PostFormValue("teacherId") + ".json"
		err = ioutil.WriteFile(writePath, output, 0644)
		if err != nil {
			panic(err)
		}

		//设置两个cookie包含教师id和name的信息
		cookieId := http.Cookie{
			Name:     "teacherId",
			Value:    r.PostFormValue("teacherId"),
			HttpOnly: false,
		}
		http.SetCookie(w, &cookieId)

		cookieName := http.Cookie{
			Name:     "teacherName",
			Value:    teacherName,
			HttpOnly: false,
		}
		http.SetCookie(w, &cookieName)

		t.Execute(w, teacherName)
	}
}
