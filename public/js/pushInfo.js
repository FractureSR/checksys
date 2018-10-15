function pushInfo(){
	console.log("in push")
	var teacherId = getCookie("teacherId")
	var date = new Date()
	var weekday = new Array(7)
	weekday[0] = "sunday"
	weekday[1] = "monday"
	weekday[2] = "tuesday"
	weekday[3] = "wednesday"
	weekday[4] = "thursday"
	weekday[5] = "friday"
	weekday[6] = "saturday"
	
	var day = weekday[date.getDay()]
	
	var time = new Array(4)
	time[0] = "8:30"
	time[1] = "10:40"
	time[2] = "14:00"
	time[3] = "23:46"
	
	tString = date.toString()
	var t_start = tString.indexOf(date.getFullYear()) + 5
	var t_end = tString.indexOf(":") + 3
	var t = tString.substring(t_start,t_end)
	//以上过程与viewLecture（）全部一样
	for(var i=0;i<4;i++) {
		//寻找时间对应的课程
		if(inTenMin(time[i],t)) {
			//生成一个结构
			var obj = {
				"belonging":teacherId,
				"day":day,
				"sequence":i,
				"lat" : Lat,
				"lng" : Lng,
				"accu": Accu
			}
			console.log(obj)
			//虚拟表单技术
			var virtualForm = document.createElement("form")
			//提交到/check
			virtualForm.action = "/check"
			virtualForm.method = "post"
			virtualForm.style.display = "none"
			
			for(var x in obj) {
				var row = document.createElement("textarea")
				row.name = x;
				row.value = obj[x]
				virtualForm.appendChild(row)
			}
			
			document.body.appendChild(virtualForm)
			virtualForm.submit()
			
			return virtualForm
		}
	}
}