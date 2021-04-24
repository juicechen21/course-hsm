package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-micro/registry/consul"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry"
	go_micro_srv_uavdata "hsm/webservice/proto/uavdata"
	go_micro_srv_user "hsm/webservice/proto/user"
	"hsm/webservice/util"
	"net/http"
	"strconv"
)

var service micro.Service

func init()  {
	// 初始化服务发现 consul
	consulReg := consul.NewRegistry(registry.Addrs("192.168.5.88:8500"))
	// 初始化micro服务对象，指定consul为服务发现
	service = micro.NewService(
		micro.Registry(consulReg),
	)
}

// LonginHandler 用户登录
func LonginHandler(c *gin.Context)  {
	username := c.PostForm("username")
	password := c.PostForm("password")

	//// 初始化服务发现 consul
	//consulReg := consul.NewRegistry(registry.Addrs("192.168.5.88:8500"))
	//
	//// 初始化micro服务对象，指定consul为服务发现
	//service := micro.NewService(
	//	micro.Registry(consulReg),
	//)

	// 初始化客户端
	microClient := go_micro_srv_user.NewUserService("go.micro.srv.user",service.Client())

	// 远程调用服务
	resp, err := microClient.Login(context.TODO(), &go_micro_srv_user.Request{
		Username: username,
		Password: password,
	})

	if err != nil {
		fmt.Println("call err: ", err)
	}
	c.JSON(http.StatusOK,gin.H{
		"status":"",
		"message":"ok",
		"return":resp,
	})
}

// PushUavDataHandler 上传无人机数据
func PushUavDataHandler(c *gin.Context) {

	bodyData, err := c.GetRawData() // 从c.Request.Body读取请求数据
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
		})
		return
	}
	resq,err := GetFormatUavData(bodyData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  http.StatusInternalServerError,
			"message": err,
		})
		return
	}

	// 初始化客户端
	microClient := go_micro_srv_uavdata.NewUavdataService("go.micro.srv.uavdata",service.Client())
	// 远程调用服务
	resp,err := microClient.Call(context.TODO(),resq)

	if err != nil {
		c.JSON(http.StatusInternalServerError,gin.H{
			"status":http.StatusInternalServerError,
			"message":err,
		})
		return
	}
	var m map[string]interface{}
	json.Unmarshal([]byte(resp.Msg),&m)
	c.JSON(http.StatusOK,gin.H{
		"status":http.StatusOK,
		"message":"ok",
		"return":m,
	})
}


// SaveUavData 保存无人机数据 落盘 mysql 存储
func GetFormatUavData(bodyData []byte) (*go_micro_srv_uavdata.Request,error){
	// 定义map或结构体
	var uavMapData map[string]interface{}
	// 反序列化
	err := json.Unmarshal(bodyData, &uavMapData)
	if err != nil {
		return nil,err
	}

	var (
		systemCodeStr       string
		deviceHardIdStr 	string
		deviceTypeInt       int64
		airSpeedFloat       float32
		altitudeFloat       float32
		barometerTempFloat  float32 //气压计温度
		battaryRemainFloat  float32 //电池剩余电量（%）
		climbRateFloat      float32 //爬升率
		currentFloat        float32 //电流
		dateTimeInt         int64   //飞行时间
		distanceToHomeFloat float32 //到Home点距离
		distanceToNextFloat float32 //到下一任务点距离
		flightDistanceFloat float32 //已飞公里数（单位KM）注意客户端上报过来的单位是什么
		flightModeStr       string  //飞行模式
		flightSortieStr     string  //飞行架次
		flightStateInt      int64   //飞行状态
		flightTimeFloat     float32 //已飞行时间（单位：min）
		groundSpeedFloat    float32 //对地速度
		heightFloat         float32 //高度
		imuTempFloat        float32 //处理器温度
		isLocationInt       int64   //是否定位成功
		latitudeStr         string  //纬度
		longitudeStr        string  //经度
		pitchFloat          float32 //俯仰
		rollFloat           float32 //横滚
		satCountFloat       float32 //卫星数
		unmannedIdInt       int64   //无人机类型
		voltageStr          string  //电池电压
		yawStr              string  //偏航
		mountTypeStr        string  //挂载类型
		uidInt              int64   //用户ID
		taskIdInt           int64   //任务ID
		platformTypeStr     string  //平台类型
		customDataStr       string  //自定义（源数据）
		currentMountTypeStr string  //挂载类型（源数据）
		mountInfoStr        string  //挂载数据（源数据）
		videoInfoStr        string  //视频数据（源数据）
		cmdGasListStr       string  //气体数据（源数据）
		dataStr             string  //上报数据（源数据）
		time                int64   = util.GetTimeNowUnix()
	)

	// 系统编码
	systemCode, ok := uavMapData["systemCode"]
	if ok {
		systemCodeStr = util.Strval(systemCode)
	} else {
		systemCodeStr = "MMC"
	}

	// 设备编号
	deviceHardId, ok := uavMapData["deviceHardId"]
	if ok {
		deviceHardIdStr = util.Strval(deviceHardId)
	} else {
		deviceHardIdStr = ""
	}

	// 设备类型
	deviceType, ok := uavMapData["deviceType"]
	if ok {
		deviceTypeStr := util.Strval(deviceType)
		deviceTypeInt, _ = strconv.ParseInt(deviceTypeStr, 10, 64) // 字符串 转 int 类型
	} else {
		deviceTypeInt = 1
	}

	// 上报源数据
	data, ok := uavMapData["data"]
	if ok {
		// dataStr   string  //上报数据（源数据）
		dataStr = util.Strval(data)
		uavInfo, ok := data.(map[string]interface{})["uavInfo"]
		if ok {
			// airSpeedFloat 空气速度
			airSpeed, ok := uavInfo.(map[string]interface{})["airSpeed"]
			if ok {
				airSpeedFloats, _ := strconv.ParseFloat(util.Strval(airSpeed), 64) // 字符串 转 浮点类型
				airSpeedFloat = float32(airSpeedFloats)
			} else {
				airSpeedFloat = 0.00
			}
			// altitude 海拔
			altitude, ok := uavInfo.(map[string]interface{})["altitude"]
			if ok {
				altitudeFloats, _ := strconv.ParseFloat(util.Strval(altitude), 64)
				altitudeFloat = float32(altitudeFloats)
			} else {
				altitudeFloat = 0.00
			}
			// barometerTempFloat float64 //气压计温度
			barometerTemp, ok := uavInfo.(map[string]interface{})["barometerTemp"]
			if ok {
				barometerTempFloats, _ := strconv.ParseFloat(util.Strval(barometerTemp), 64)
				barometerTempFloat = float32(barometerTempFloats)
			} else {
				barometerTempFloat = 0.00
			}
			//电池剩余电量（%）
			battaryRemain, ok := uavInfo.(map[string]interface{})["battaryRemain"]
			if ok {
				battaryRemainFloats, _ := strconv.ParseFloat(util.Strval(battaryRemain), 64)
				battaryRemainFloat = float32(battaryRemainFloats)
			} else {
				battaryRemainFloat = 0.00
			}
			// climbRateFloat          float64 //爬升率
			climbRate, ok := uavInfo.(map[string]interface{})["climbRate"]
			if ok {
				climbRateFloats, _ := strconv.ParseFloat(util.Strval(climbRate), 64)
				climbRateFloat = float32(climbRateFloats)
			} else {
				climbRateFloat = 0.00
			}
			// currentFloat            float64 //电流
			current, ok := uavInfo.(map[string]interface{})["current"]
			if ok {
				currentFloats, _ := strconv.ParseFloat(util.Strval(current), 64)
				currentFloat = float32(currentFloats)
			} else {
				currentFloat = 0.00
			}
			// dateTimeInt  int64   //飞行时间
			dateTime, ok := uavInfo.(map[string]interface{})["dateTime"]
			if ok {
				dateTimeStr := util.Strval(dateTime)
				dateTimeInt, _ = strconv.ParseInt(dateTimeStr, 10, 64)
			} else {
				dateTimeInt = 0
			}
			// distanceToHomeFloat float64 //到Home点距离
			distanceToHome, ok := uavInfo.(map[string]interface{})["distanceToHome"]
			if ok {
				distanceToHomeFloats, _ := strconv.ParseFloat(util.Strval(distanceToHome), 64)
				distanceToHomeFloat = float32(distanceToHomeFloats)
			} else {
				distanceToHomeFloat = 0.00
			}
			// distanceToNextFloat float64 //到下一任务点距离
			distanceToNext, ok := uavInfo.(map[string]interface{})["distanceToNext"]
			if ok {
				distanceToNextFloats, _ := strconv.ParseFloat(util.Strval(distanceToNext), 64)
				distanceToNextFloat = float32(distanceToNextFloats)
			} else {
				distanceToNextFloat = 0.00
			}
			// flightDistanceFloat float64 //已飞公里数（单位KM）注意客户端上报过来的单位是什么
			flightDistance, ok := uavInfo.(map[string]interface{})["flightDistance"]
			if ok {
				flightDistanceFloats, _ := strconv.ParseFloat(util.Strval(flightDistance), 64)
				flightDistanceFloat = float32(flightDistanceFloats)
			} else {
				flightDistanceFloat = 0.00
			}
			// flightModeStr     string  //飞行模式
			flightMode, ok := uavInfo.(map[string]interface{})["flightMode"]
			if ok {
				flightModeStr = util.Strval(flightMode)
			} else {
				flightModeStr = "0"
			}
			// flightSortieStr   string  //飞行架次
			flightSortie, ok := uavInfo.(map[string]interface{})["flightSortie"]
			if ok {
				flightSortieStr = util.Strval(flightSortie)
			} else {
				flightSortieStr = "0"
			}
			// flightStateInt    int64   //飞行状态
			flightState, ok := uavInfo.(map[string]interface{})["flightState"]
			if ok {
				flightStateStr := util.Strval(flightState)
				flightStateInt, _ = strconv.ParseInt(flightStateStr, 10, 64)
			} else {
				flightStateInt = 0
			}
			// flightTimeFloat     float64 //已飞行时间（单位：min）
			flightTime, ok := uavInfo.(map[string]interface{})["flightTime"]
			if ok {
				flightTimeFloats, _ := strconv.ParseFloat(util.Strval(flightTime), 64)
				flightTimeFloat = float32(flightTimeFloats)
			} else {
				flightTimeFloat = 0
			}
			// groundSpeedFloat    float64 //对地速度
			groundSpeed, ok := uavInfo.(map[string]interface{})["groundSpeed"]
			if ok {
				groundSpeedFloats, _ := strconv.ParseFloat(util.Strval(groundSpeed), 64)
				groundSpeedFloat = float32(groundSpeedFloats)
			} else {
				groundSpeedFloat = 0.00
			}
			// heightFloat         float64 //高度
			height, ok := uavInfo.(map[string]interface{})["height"]
			if ok {
				heightFloats, _ := strconv.ParseFloat(util.Strval(height), 64)
				heightFloat = float32(heightFloats)
			} else {
				heightFloat = 0.00
			}
			// imuTempFloat        float64 //处理器温度
			imuTemp, ok := uavInfo.(map[string]interface{})["imuTemp"]
			if ok {
				imuTempFloats, _ := strconv.ParseFloat(util.Strval(imuTemp), 64)
				imuTempFloat = float32(imuTempFloats)
			} else {
				imuTempFloat = 0.00
			}
			// isLocationInt     int64   //是否定位成功
			isLocation, ok := uavInfo.(map[string]interface{})["isLocation"]
			if ok {
				isLocationStr := util.Strval(isLocation)
				isLocationInt, _ = strconv.ParseInt(isLocationStr, 10, 64)
			} else {
				isLocationInt = 0
			}
			// latitudeStr       string  //纬度
			latitude, ok := uavInfo.(map[string]interface{})["latitude"]
			if ok {
				latitudeStr = util.Strval(latitude)
			} else {
				latitudeStr = "0"
			}
			// longitudeStr      string  //经度
			longitude, ok := uavInfo.(map[string]interface{})["longitude"]
			if ok {
				longitudeStr = util.Strval(longitude)
			} else {
				longitudeStr = "0"
			}
			// pitchFloat          float64 //俯仰
			pitch, ok := uavInfo.(map[string]interface{})["pitch"]
			if ok {
				pitchFloats, _ := strconv.ParseFloat(util.Strval(pitch), 64)
				pitchFloat = float32(pitchFloats)
			} else {
				pitchFloat = 0
			}
			// rollFloat           float64 //横滚
			roll, ok := uavInfo.(map[string]interface{})["roll"]
			if ok {
				rollFloats, _ := strconv.ParseFloat(util.Strval(roll), 64)
				rollFloat = float32(rollFloats)
			} else {
				rollFloat = 0
			}
			// satCountFloat       float64 //卫星数
			satCount, ok := uavInfo.(map[string]interface{})["satCount"]
			if ok {
				satCountFloats, _ := strconv.ParseFloat(util.Strval(satCount), 64)
				satCountFloat = float32(satCountFloats)
			} else {
				satCountFloat = 0
			}
			// unmannedIdInt     int64   //无人机类型
			unmanned, ok := uavInfo.(map[string]interface{})["unmanned"]
			if ok {
				unmannedStr := util.Strval(unmanned)
				unmannedIdInt, _ = strconv.ParseInt(unmannedStr, 10, 64)
			} else {
				unmannedIdInt = 0
			}
			// voltageStr        string  //电池电压
			voltage, ok := uavInfo.(map[string]interface{})["voltage"]
			if ok {
				voltageStr = util.Strval(voltage)
			} else {
				voltageStr = "0"
			}
			// yawStr            string  //偏航
			yaw, ok := uavInfo.(map[string]interface{})["yaw"]
			if ok {
				yawStr = util.Strval(yaw)
			} else {
				yawStr = "0"
			}
			// mountTypeStr      string  //挂载类型
			mountType, ok := uavInfo.(map[string]interface{})["mountType"]
			if ok {
				mountTypeStr = util.Strval(mountType)
			} else {
				mountTypeStr = "0"
			}
			// uidInt            int64   //用户ID
			uid, ok := uavInfo.(map[string]interface{})["uid"]
			if ok {
				uidStr := util.Strval(uid)
				uidInt, _ = strconv.ParseInt(uidStr, 10, 64)
			} else {
				uidInt = 0
			}
			// taskIdInt         int64   //任务ID
			taskId, ok := uavInfo.(map[string]interface{})["taskId"]
			if ok {
				taskIdStr := util.Strval(taskId)
				taskIdInt, _ = strconv.ParseInt(taskIdStr, 10, 64)
			} else {
				taskIdInt = 0
			}
			// platformTypeStr   string  //平台类型
			platformType, ok := uavInfo.(map[string]interface{})["platformType"]
			if ok {
				platformTypeStr = util.Strval(platformType)
			} else {
				platformTypeStr = "0"
			}
			// customDataStr       string  //自定义（源数据）
			customData, ok := uavInfo.(map[string]interface{})["customData"]
			if ok {
				customDataStr = util.Strval(customData)
			} else {
				customDataStr = "[]"
			}
			// currentMountTypeStr string  //挂载类型（源数据）
			currentMountType, ok := uavInfo.(map[string]interface{})["currentMountType"]
			if ok {
				currentMountTypeStr = util.Strval(currentMountType)
			} else {
				currentMountTypeStr = "[]"
			}
		}
	}
	// videoInfoStr        string  //视频数据（源数据）
	videoInfo, ok := data.(map[string]interface{})["videoInfo"]
	if ok {
		videoInfoStr = util.Strval(videoInfo)
	} else {
		videoInfoStr = "[]"
	}
	// cmdGasListStr       string  //气体数据（源数据）  mountInfoStr        string  //挂载数据（源数据）
	cmdGasList, ok := data.(map[string]interface{})["mountInfo"]
	if ok {
		cmdGasListStr = util.Strval(cmdGasList)
		mountInfoStr = cmdGasListStr
	} else {
		mountInfoStr, cmdGasListStr = "[]", "[]"
	}

	uavOneData := &go_micro_srv_uavdata.Request{
			SystemCode:       systemCodeStr,
			DeviceType:       deviceTypeInt,
			DeviceHardId:     deviceHardIdStr,
			AirSpeed:         airSpeedFloat,
			Altitude:         altitudeFloat,
			BarometerTemp:    barometerTempFloat,
			BattaryRemain:    battaryRemainFloat,
			ClimbRate:        climbRateFloat,
			Current:          currentFloat,
			DateTime:         dateTimeInt,
			DistanceToHome:   distanceToHomeFloat,
			DistanceToNext:   distanceToNextFloat,
			FlightDistance:   flightDistanceFloat,
			FlightMode:       flightModeStr,
			FlightSortie:     flightSortieStr,
			FlightState:      flightStateInt,
			FlightTime:       flightTimeFloat,
			GroundSpeed:      groundSpeedFloat,
			Height:           heightFloat,
			ImuTemp:          imuTempFloat,
			IsLocation:       isLocationInt,
			Latitude:         latitudeStr,
			Longitude:        longitudeStr,
			Pitch:            pitchFloat,
			Roll:             rollFloat,
			SatCount:         satCountFloat,
			UnmannedId:       unmannedIdInt,
			Voltage:          voltageStr,
			Yaw:              yawStr,
			MountType:        mountTypeStr,
			Uid:              uidInt,
			TaskId:           taskIdInt,
			PlatformType:     platformTypeStr,
			CustomData:       customDataStr,
			CurrentMountType: currentMountTypeStr,
			MountInfo:        mountInfoStr,
			VideoInfo:        videoInfoStr,
			CmdGasList:       cmdGasListStr,
			Data:             dataStr,
			Time:             time,
		}
	return uavOneData,nil
}




