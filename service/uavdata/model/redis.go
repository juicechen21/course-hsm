package model

import "context"

type PushUavData struct {
	SystemCode       string  `db:"systemCode"`       //系统编码
	DeviceType       int64   `db:"deviceType"`       //设备类型：1无人机 2地面站 3挂载
	DeviceHardId     string  `db:"deviceHardId"`     //设备id
	AirSpeed         float64 `db:"airSpeed"`         //空气速度（单位：m/s）
	Altitude         float64 `db:"altitude"`         //海拔
	BarometerTemp    float64 `db:"barometerTemp"`    //气压计温度
	BattaryRemain    float64 `db:"battaryRemain"`    //电池剩余电量（%）
	ClimbRate        float64 `db:"climbRate"`        //爬升率
	Current          float64 `db:"current"`          //电流
	DateTime         int64   `db:"dateTime"`         //飞行时间
	DistanceToHome   float64 `db:"distanceToHome"`   //到Home点距离
	DistanceToNext   float64 `db:"distanceToNext"`   //到下一任务点距离
	FlightDistance   float64 `db:"flightDistance"`   //已飞公里数（单位KM）注意客户端上报过来的单位是什么
	FlightMode       string  `db:"flightMode"`       //飞行模式
	FlightSortie     string  `db:"flightSortie"`     //飞行架次
	FlightState      int64   `db:"flightState"`      //飞行状态
	FlightTime       float64 `db:"flightTime"`       //已飞行时间（单位：min）
	GroundSpeed      float64 `db:"groundSpeed"`      //对地速度
	Height           float64 `db:"height"`           //高度
	ImuTemp          float64 `db:"imuTemp"`          //处理器温度
	IsLocation       int64   `db:"isLocation"`       //是否定位成功
	Latitude         string  `db:"latitude"`         //纬度
	Longitude        string  `db:"longitude"`        //经度
	Pitch            float64 `db:"pitch"`            //俯仰
	Roll             float64 `db:"roll"`             //横滚
	SatCount         float64 `db:"satCount"`         //卫星数
	UnmannedId       int64   `db:"unmannedId"`       //无人机类型
	Voltage          string  `db:"voltage"`          //电池电压
	Yaw              string  `db:"yaw"`              //偏航
	MountType        string  `db:"mountType"`        //挂载类型
	Uid              int64   `db:"uid"`              //用户ID
	TaskId           int64   `db:"taskId"`           //任务ID
	PlatformType     string  `db:"platformType"`     //平台类型
	CustomData       string  `db:"customData"`       //自定义（源数据）
	CurrentMountType string  `db:"currentMountType"` //挂载类型（源数据）
	MountInfo        string  `db:"mountInfo"`        //挂载数据（源数据）
	VideoInfo        string  `db:"videoInfo"`        //视频数据（源数据）
	CmdGasList       string  `db:"cmdGasList"`       //气体数据（源数据）
	Data             string  `db:"data"`             //上报数据（源数据）
	Time             int64   `db:"time"`
}


// LPush 将数据加入无序集合队列
func LPushData(key string,data interface{}) (val int64) {
	ctx := context.Background()
	val = redisdb.LPush(ctx,key,data).Val()
	return
}
