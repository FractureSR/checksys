package threads

import (
	"fmt"
	"sync"

	"checksys2.0/apratus"
)

//原子锁和线程计数
var mutex sync.Mutex
var CountThreads = 0

//用一个可以通过pushInfo的chan和外界通信
func Core(pushInfoDeliver chan apratus.PushInfo) {
	//core启动时第一个线程启动了，CountThreads被多个线程共用
	mutex.Lock()
	{
		CountThreads++
	}
	mutex.Unlock()
	fmt.Println("Core started")
	for {
		//准备接收pushInfo信息
		getPushInfo, ok := <-pushInfoDeliver
		if !ok {
			fmt.Println("Core shutdown")
		}
		fmt.Println("Info get")
		fmt.Println(getPushInfo)
		//得到了pushInfo的信息现在要把信息存入资源池
		mutex.Lock()
		{
			//以教师的唯一id为标识将推送信息存入映射并且新建一个记录器存入记录器映射，最后为这个推送启动一个线程，这个线程与记录器之间用记录器的通道作为绑定标识
			apratus.LectureInfo[getPushInfo.Id] = getPushInfo
			apratus.Recorders[getPushInfo.Id] = apratus.NewRecorder()
			go FuncRecorder(apratus.Recorders[getPushInfo.Id].PrivateChan, getPushInfo.Id)
		}
		mutex.Unlock()
	}
}

func ViceCore(restudentInfoDeliver chan apratus.RestudentInfo) {
	//线程+1
	mutex.Lock()
	{
		CountThreads++
	}
	mutex.Unlock()
	fmt.Println("ViceCore started")
	for {
		//等待学生返回信息
		getStudentInfo, ok := <-restudentInfoDeliver
		if !ok {
			fmt.Println("ViceCore shutdown")
		}
		fmt.Println("StudentInfo get")
		fmt.Println(getStudentInfo)
		//从记录器中找到学生选择的老师的课程，将信息发送到对应的通道中
		mutex.Lock()
		{
			apratus.Recorders[getStudentInfo.TeacherId].PrivateChan <- getStudentInfo
		}
		mutex.Unlock()
	}
}
