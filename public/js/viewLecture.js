var Lat
var Lng
var Accu

//这个函数处理不了23点到0点的情况
function inTenMin(t1,t2){
	//将t1,t2分解成两个小时，两个分钟
	var h1 = parseInt(t1.substring(0,t1.indexOf(":")))
	var m1 = parseInt(t1.substring(t1.indexOf(":")+1))
	var h2 = parseInt(t2.substring(0,t2.indexOf(":")))
	var m2 = parseInt(t2.substring(t2.indexOf(":")+1))
	if(h1 - h2 == 1) {
		if(60+m1-m2 <= 10) return true
		else return false
	}
	else if(h1 == h2) {
		if(Math.abs(m1-m2) <= 10) return true
		else return false
	}
	else if(h2 - h1 == 1) {
		if(60+m2-m1 <= 10) return true
		else return false
	}
	else return false
}

function viewLecture(){
	var teacherId = getCookie("teacherId")
	var route = "/static/temp/" + teacherId + ".json"
	//取得静态课表
	$.get(route,function (data) {
		//解码
		var obj = eval("(" + data + ")")
		var date = new Date()
		//设置数组weekday
		var weekday = new Array(7)
		weekday[0] = "sunday"
		weekday[1] = "monday"
		weekday[2] = "tuesday"
		weekday[3] = "wednesday"
		weekday[4] = "thursday"
		weekday[5] = "friday"
		weekday[6] = "saturday"
		
		//设置数组上课时间
		var time = new Array(4)
		time[0] = "8:30"
		time[1] = "10:40"
		time[2] = "14:00"
		time[3] = "23:46"
		
		var check = true
		for(day in obj) {
			//首先找到和今天一样的weekday
			if(day == weekday[date.getDay()]) {
				//获取具体时间并字符串化
				tString = date.toString()
				//start在'年'之后的第五位
				var t_start = tString.indexOf(date.getFullYear()) + 5
				//end是找到时间中间的：后的第三位
				var t_end = tString.indexOf(":") + 3
				//t为截取的时间格式为x(x):xx，:前面可以是一位或者两位
				var t = tString.substring(t_start,t_end)
				for(var i=0;i<4;i++) {
					//用time数组中的四个之间逐个与截取的时间比较，是否与指定时间差十分钟
					if(inTenMin(time[i],t)) {
						if(obj[day][i] == "&nbsp") break //找到了时间符合但是这一节是空课也要排除
						check = false
						var msg = "<tr>" + "<th>LECTURE</th>" + "<th>TEACHER</th>" + "<th>AUDIENCE</th>" +"<th>STATUS</th>" + "<th>PUSH</th>" + "</tr>" 
						var lecture = obj[day][i].substring(0,obj[day][i].indexOf("<br />"))
						var audience = obj[day][i].substring(obj[day][i].indexOf(">")+1)
						msg = msg + "<tr>" + "<td>" + lecture + "</td>" + "<td>" + getCookie("teacherName") +"</td>" + "<td>" + audience +"</td>" +"<td>ready</td>" + "<td><button type=\"button\" onclick=\"pushInfo()\" style=\"height:30px;width:80px\">PUSH</button></td>" +"</tr>"
						//替换成发布信息，包含一个按钮，启动pushInfo（）
						document.getElementById("available").innerHTML=msg
					}
				}
				break
			}
		}
		//未找到时间对应或者没有课
		if(check) document.getElementById("available").innerHTML="<div style=\"text-align:center\">NO LECTURE IS BEGINING AT THE MOMENT</div>"
	})
}

function getPos(){
	console.log("in getPos")
	//启动封装好的腾讯定位接口，设置超时
	var geolocation = new qq.maps.Geolocation("2V2BZ-2S366-LAJSR-EI25A-C22O7-PCBSW","checksys")
	var options = {timeout: 8000}
	//有两个回调，一个成功，一个失败
	geolocation.getLocation(locateSuccess,locateFail,options)
}

function locateSuccess(position){
	//定位成功时但是精度太低会提示
	if(position.accuracy > 500){
		alert("TERRIBLE GPS ACCURACY")
	}
	//position是一个结构，将他json化
	Lat = JSON.stringify(position.lat)
	Lng = JSON.stringify(position.lng)
	Accu = JSON.stringify(position.accuracy)
}

function locateFail(){
	//定位失败时提示定位失败并试图关闭窗口
	alert("Locate Failed.")
	window.opener = null;
	window.open("","_self","")
	window.close()
}