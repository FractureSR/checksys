package threads

import (
	"fmt"
	"time"

	"checksys2.0/apratus"
)

func FuncRecorder(data chan apratus.RestudentInfo, serve string) {
	//线程增加
	mutex.Lock()
	{
		CountThreads++
	}
	mutex.Unlock()
	fmt.Println("New threads started!\nCurrent Threads:", CountThreads, "Two for cores.")
	fmt.Println("Working for", serve)
	//为自己启动一个时间监视器
	go TimeWatcher(serve, 5)
	for {
		//等待vicecore传来数据
		getData, ok := <-data
		//一旦通道被关闭整个记录器处理函数会进入收尾工作并结束
		if !ok {
			fmt.Println("Info transfer mistake.Serving,", serve)
			break
		}
		fmt.Println(getData)
		//将得到的数据存入记录器的数据区
		mutex.Lock()
		{
			apratus.Recorders[serve].Append(getData)
			fmt.Println("Add success!", apratus.Recorders[serve].GetInfo())
		}
		mutex.Unlock()
	}
	mutex.Lock()
	{
		CountThreads--                     //线程减少
		delete(apratus.LectureInfo, serve) //从课程映射里删除记录
		fmt.Println("LectureInfo left", apratus.LectureInfo)
		delete(apratus.Recorders, serve) //从记录器映射里删除记录器数据
		fmt.Println("Recorders left", apratus.Recorders)
	}
	mutex.Unlock()
	fmt.Println("Timeout,work done,threads fall out.Current threads:", CountThreads, "Two for cores.")
	fmt.Println("Woring finished,", serve)
}

//时间监视器,在经过waitTime后关闭打开自己的FuncRecorder绑定的通道
func TimeWatcher(serve string, waitTime time.Duration) {
	time.Sleep(waitTime * time.Minute)
	fmt.Println("work done!")
	mutex.Lock()
	{
		close(apratus.Recorders[serve].PrivateChan)
	}
	mutex.Unlock()
}
