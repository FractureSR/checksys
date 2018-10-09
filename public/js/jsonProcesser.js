// cookie的数据以键值对间用;分开的办法存储在同一个文本文档中
function getCookie(c_name){
	//首先要求存储cookie的文档有数据
	if (document.cookie.length>0) {
		//查找名字为c_name的cookie，首先查找c_name + "=" 的字符串
		c_start=document.cookie.indexOf(c_name + "=")
		//是否找到了这个名称的cookie
		if (c_start != -1) {
			//接下来取键值对"="后的部分，调整start到等号后
			c_start=c_start + c_name.length + 1
			//查找c_start开始后第一个;即为结尾
			c_end=document.cookie.indexOf(";",c_start)
			//找不到分号说明该cookie在文档尾，此时尾部就是整个文档的尾
			if (c_end==-1) c_end=document.cookie.length
			//找到字串，因为cookie进行过编码所以用unescape进行解码并返回
			return unescape(document.cookie.substring(c_start,c_end))
		}
	}
	//文档为空或者找不到指定名字的cookie时返回空字符串
	return ""
}

function jsonProcesser(){
	console.log("powered")
	//通过teacherId的cookie组合出服务器中课程表的静态存放地址
	var teacherId = getCookie("teacherId")
	var route = "/static/temp/" + teacherId + ".json"
	//jquery的get方法获取课程表的值并放在data中
	$.get(route,function (data) {
		console.log(data)
		//eval只能解码在（）间的内容，把json填充的结构中
		var obj = eval("(" + data + ")")
		console.log(obj)
		var table = ""
		//表头
		table += "<tr><th></th><th>Sunday</th><th>Monday</th><th>Tuesday</th><th>Wednesday</th><th>Thursday</th><th>Friday</th><th>Saturday</th></tr>"
		//午休前
		for(var i=0;i<2;i++) {
			table += "<tr>"
			table = table + "<th>" + (i+1).toString() + "</th>"
			for (day in obj) {
				table += "<td>" + obj[day][i] + "</td>"
			}
			table += "</tr>"
		}
		//午休
		table += "<tr><td colspan=\"8\" align=\"center\">lunch time</td></tr>"
		//午休后
		for(var i=0;i<2;i++) {
			table += "<tr>"
			table = table + "<th>" + (i+3).toString() + "</th>"
			for (day in obj) {
				table += "<td>" + obj[day][i+2] + "</td>"
			}
			table += "</tr>"
		}
		//填充
		document.getElementById("table").innerHTML=table
	})
}