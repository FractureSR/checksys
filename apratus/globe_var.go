package apratus

//课表，每个weekday是一个字符串切片，包含四节课
type Timetable struct {
	Sunday    []string `json:"sunday"`
	Monday    []string `json:"monday"`
	Tuesday   []string `json:"tuesday"`
	Wednesday []string `json:"wednesday"`
	Thursday  []string `json:"thursday"`
	Friday    []string `json:"friday"`
	Saturday  []string `json:"saturday"`
}

//仿照定位信息声明的定位信息结构
type Locate struct {
	Module   string  `json:"module"`
	Type     string  `json:"type"`
	Adcode   string  `json:"adcode"`
	Nation   string  `json:"nation"`
	Province string  `json:"province"`
	City     string  `json:"city"`
	District string  `json:"district"`
	Addr     string  `json:"addr"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	Accuracy int     `json:"accuracy"`
}

type PushInfo struct {
	Lecture  string
	Teacher  string
	Audience string
	Id       string
	Lat      string
	Lng      string
	Accu     string
}

type StudentInfo struct {
	Name    string `json:"Name"`
	College string `json:"College"`
	Major   string `json:"Major"`
	Class   string `json:"Class"`
	Id      string `json:"Id"`
}

type RestudentInfo struct {
	Basic     StudentInfo `json:"basic"`
	TeacherId string      `json:"teacherId"`
	//SeatNum   string      `json:"seatNum"`
	//RoomType  string      `json:"roomType"`
	SLat  string `json:"sLat"`
	SLng  string `json:"sLng"`
	SAccu string `json:"sAccu"`
}

//记录器结构，包含一个通信用的通道和一个用于记录学生返回的数据的切片
type Recorder struct {
	PrivateChan chan RestudentInfo
	CheckInData []RestudentInfo
}

type Coordinate struct {
	Lat      float64
	Lng      float64
	Accuracy int
}

type DistInfo struct {
	Distance   float64
	Accu_level string
}

//返回一个Recoder结构，构造好了通道并把数据区清零
func NewRecorder() Recorder {
	return Recorder{
		PrivateChan: make(chan RestudentInfo),
		CheckInData: nil,
	}
}

func (r Recorder) Append(restudentInfo RestudentInfo) {
	r.CheckInData = append(r.CheckInData, restudentInfo)
}

//方法，返回数据区数据
func (r Recorder) GetInfo() []RestudentInfo {
	return r.CheckInData
}

var PushInfoDeliver chan PushInfo

var RestudentInfoDeliver chan RestudentInfo

var LectureInfo map[string]PushInfo

var Recorders map[string]Recorder
