function startTime()
			{
			var today=new Date()
			var y=today.getFullYear()
			var mon=today.getMonth() + 1
			var d=today.getDay()+1
			var h=today.getHours()
			var min=today.getMinutes()
			var s=today.getSeconds()
			// add a zero in front of numbers<10
			min=checkTime(min)
			s=checkTime(s)
			document.getElementById('timedisplay').innerHTML=y+"年"+mon+"月"+d+"日"+h+":"+min+":"+s
			t=setTimeout('startTime()',500)
			}

function checkTime(i)
			{
			if (i<10) 
			{i="0" + i}
			return i
			}