var slat
var slng
var saccu
var TeacherId

function getCookie(c_name){
	if (document.cookie.length>0) {
		c_start=document.cookie.indexOf(c_name + "=")
		if (c_start != -1) {
			c_start=c_start + c_name.length + 1
			c_end=document.cookie.indexOf(";",c_start)
			if (c_end==-1) c_end=document.cookie.length
			return unescape(document.cookie.substring(c_start,c_end))
		}
	}
	return ""
}

//jquery的post方法从后端处理器函数中拉取数据，并生成欢迎语句
function getStuInfo() {
	//通过学生id访问/studentInfo,结果存入response
	$.post("/studentInfo", {id : getCookie("stuId")}, function(response, status) {
		document.getElementById("stuInfo").innerHTML="<pre>" + "	你好，" + response.Name + "	学院：" + response.College + "	专业：" + response.Major + "	  班级：" + response.Class + "	学号：" + response.Id +"</pre>"
	}, "json")
}

function getLectureInfo() {
	console.log("in get LectureInfo")
	//get方法访问/lectureInfo
	$.get("/lectureInfo", function (response) {
		console.log(response)
		if(response == null) {
			var msg = "<h2><strong>THERE IS NO LECTURE BEGINING.</strong></h2>"
			document.getElementById("lectureInfo").innerHTML = msg
		}
		//构造课程表格，提供CheckIn()按钮
		else {
			console.log(response)
			var msg = "<tr>" + "<th>No.</th>" + "<th>LECTURE</th>" + "<th>TEACHER</th>" + "<th>AUDIENCE</th>" + "<th>CHECK</th>" + "</tr>" 
			for(var i = 0; i < response.length; i++) {
				msg = msg + "<tr>" + "<td>" + (i+1).toString() + "</td>" + "<td>" + response[i].Lecture + "</td>" + "<td>" + response[i].Teacher + "</td>" + "<td>" + response[i].Audience + "</td>" + "<td>" + "<button id=\"" + response[i].Id + "\"onclick=\"getId(this.id),CheckIn()\" style=\"height:30px;width:80px;\">CHECK</button>" + "</td>" + "</tr>"
			}
			document.getElementById("lectureInfo").innerHTML = msg
		}
	})
}

function CheckIn() {
	//post方法访问/studnetInfo
	$.post("/studentInfo", {id : getCookie("stuId")}, function(response, status) {
		var obj = {
			"basic" : response,
			"teacherId" : TeacherId,
			"sLat" : slat,
			"sLng" : slng,
			"sAccu": saccu
		}
		//构建成返回学生信息并再次post访问，/studentCheck
		match = JSON.stringify(obj)
		console.log(match)
		$.post("/studentCheck", {data : match}, function(response, status) {
			console.log(response)
			if(response == "Major Error" || response == "Class Error") alert("It seems you are not scheduled to attend the lecture.QAQ")
			else if(response.Distance <=150){
				alert("You have registered successfully!")
			}
			else {
				alert("There is something wrong with your positioning.Please inform the teacher personally.")
			}
		})
	}, "json")
}

function getId(id) {
	TeacherId = id
}

//与getPos相同
function getPosAlter(){
	console.log("in getPos")
	var geolocation = new qq.maps.Geolocation("2V2BZ-2S366-LAJSR-EI25A-C22O7-PCBSW","checksys")
	var options = {timeout: 8000}
	geolocation.getLocation(locateSuccessAlter,locateFailAlter,options)
}

//成功回调
function locateSuccessAlter(position){
	if(position.accuracy > 500){
		alert("BAD GPS ACCURACY")
	}
	slat = JSON.stringify(position.lat)
	slng = JSON.stringify(position.lng)
	saccu = JSON.stringify(position.accuracy)
}

//失败回调
function locateFailAlter(){
	alert("Locate Failed.")
	window.opener = null;
	window.open("","_self","")
	window.close()
}