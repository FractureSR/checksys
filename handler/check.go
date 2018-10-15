package handler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	//"html/template"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"checksys2.0/apratus"
	_ "github.com/go-sql-driver/mysql"
)

//原子锁
var mutex sync.Mutex

func Check(w http.ResponseWriter, r *http.Request) {
	var Alocate apratus.Locate
	//将地址信息解码到Alocate中
	json.Unmarshal([]byte(r.PostFormValue("location")), &Alocate)
	//打开数据库
	db, err := sql.Open("mysql", "root:Neu1923#@tcp(localhost:3306)/checksys_data?charset=utf8")
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

	var sqlSentence string
	var lectureInfo string
	//选出被push的一节课
	sqlSentence = "SELECT " + r.PostFormValue("day") + " FROM timetable WHERE belonging=\"" + r.PostFormValue("belonging") + "\" AND sequence=\"" + r.PostFormValue("sequence") + "\""
	fmt.Println(sqlSentence)
	err = db.QueryRow(sqlSentence).Scan(&lectureInfo)
	if err != nil {
		panic(err)
	}
	var lecture string
	var audience string
	//lectureInfo包含课程名称和上课班级，他们之间用一个<br>分开，现在分别取出两个信息
	lecture = apratus.SubString(lectureInfo, 0, strings.Index(lectureInfo, "<br />"))
	audience = apratus.SubString(lectureInfo, strings.Index(lectureInfo, ">")+1, len(lectureInfo))
	fmt.Println(lecture, audience)

	//取得老师的名字
	var teacherName string
	sqlSentence = "SELECT nickname FROM teacherInfo WHERE teacherId=\"" + r.PostFormValue("belonging") + "\""
	err = db.QueryRow(sqlSentence).Scan(&teacherName)
	if err != nil {
		panic(err)
	}

	//构建pushInfo结构
	ApushInfo := apratus.PushInfo{
		Lecture:  lecture,
		Teacher:  teacherName,
		Audience: audience,
		Id:       r.PostFormValue("belonging"),
		Lat:      r.PostFormValue("lat"),
		Lng:      r.PostFormValue("lng"),
		Accu:     r.PostFormValue("accu"),
	}
	fmt.Println(ApushInfo)

	var exists bool
	//检查这个老师是否已经发布过课程
	mutex.Lock()
	{
		_, exists = apratus.LectureInfo[ApushInfo.Id]
	}
	mutex.Unlock()
	if exists { //如果已经发布了
		//t, err := template.ParseFiles("templates/")
	} else { //如果没有发布的话，发布信息将通过通道输送到core处理
		apratus.PushInfoDeliver <- ApushInfo
	}
}

func StudentCheck(w http.ResponseWriter, r *http.Request) {
	//从浏览器接收到学生返回的信息
	fmt.Println("In StudentCheck")
	fmt.Println(r.PostFormValue("data"))
	raw := []byte(r.PostFormValue("data"))
	fmt.Println(string(raw))
	var ArestudentInfo apratus.RestudentInfo
	json.Unmarshal(raw, &ArestudentInfo)
	fmt.Println(ArestudentInfo)
	var audience string
	var teacherLat, teacherLng, teacherAccu string
	mutex.Lock()
	{
		audience = apratus.LectureInfo[ArestudentInfo.TeacherId].Audience
		teacherLat = apratus.LectureInfo[ArestudentInfo.TeacherId].Lat
		teacherLng = apratus.LectureInfo[ArestudentInfo.TeacherId].Lng
		teacherAccu = apratus.LectureInfo[ArestudentInfo.TeacherId].Accu
	}
	mutex.Unlock()
	//对比专业和班级
	var major_a = apratus.SubString(audience, 0, strings.Index(audience, "-")-4)
	var major_b = apratus.SubString(ArestudentInfo.Basic.Class, 0, len(ArestudentInfo.Basic.Class)-4)
	var au_start = apratus.SubString(audience, strings.Index(audience, "-")-4, strings.Index(audience, "-"))
	var au_end = apratus.SubString(audience, strings.Index(audience, "-")+1, len(audience))
	var class_num = apratus.SubString(ArestudentInfo.Basic.Class, len(ArestudentInfo.Basic.Class)-4, len(ArestudentInfo.Basic.Class))
	fmt.Println(au_start)
	fmt.Println(au_end)
	fmt.Println(class_num)
	start, err := strconv.Atoi(au_start)
	if err != nil {
		fmt.Println("strconv error")
		panic(err)
	}
	end, err := strconv.Atoi(au_end)
	if err != nil {
		fmt.Println("strconv error")
		panic(err)
	}
	value, err := strconv.Atoi(class_num)
	if err != nil {
		fmt.Println("strconv error")
		panic(err)
	}
	fmt.Println(major_a, major_b)
	if major_a == major_b {
		if value < start || value > end {
			w.Write([]byte("Class Error"))
		} else { //如果全部匹配通过坐标计算距离
			var x, _ = strconv.ParseFloat(ArestudentInfo.SLat, 64)
			var y, _ = strconv.ParseFloat(ArestudentInfo.SLng, 64)
			var z, _ = strconv.Atoi(ArestudentInfo.SAccu)
			co1 := apratus.Coordinate{
				Lat:      x,
				Lng:      y,
				Accuracy: z,
			}
			x, _ = strconv.ParseFloat(teacherLat, 64)
			y, _ = strconv.ParseFloat(teacherLng, 64)
			z, _ = strconv.Atoi(teacherAccu)
			co2 := apratus.Coordinate{
				Lat:      x,
				Lng:      y,
				Accuracy: z,
			}
			dist, accuracy := apratus.HarvenSin(co1, co2)
			var Accu_Level string
			//判断定位的可靠程度
			if accuracy > 800 {
				Accu_Level = "not reliable"
			} else if accuracy > 300 {
				Accu_Level = "normal"
			} else if accuracy > 100 {
				Accu_Level = "reliable"
			} else {
				Accu_Level = "accurate"
			}

			//构建距离信息的实例
			AdistInfo := apratus.DistInfo{
				Distance:   dist,
				Accu_level: Accu_Level,
			}

			output, err := json.MarshalIndent(&AdistInfo, "", "\t\t")
			if err != nil {
				fmt.Println("Encoding error")
				panic(err)
			}

			if AdistInfo.Distance <= 150 {
				//当满足定位条件时将准确度信息换位准确度等级
				ArestudentInfo.SAccu = AdistInfo.Accu_level
				//将信息传入副核心
				apratus.RestudentInfoDeliver <- ArestudentInfo
			}

			w.Header().Set("Content-type", "application/json")
			w.Write(output)
		}
	} else {
		w.Write([]byte("Major Error"))
	}
}
