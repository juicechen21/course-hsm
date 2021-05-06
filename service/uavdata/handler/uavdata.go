package handler

import (
	"context"
	"encoding/json"
	"fmt"

	uavdata "hsm/service/uavdata/proto/uavdata"
	"github.com/Shopify/sarama"
)

type Uavdata struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Uavdata) Call(ctx context.Context, req *uavdata.Request, rsp *uavdata.Response) error {
	v,err := json.Marshal(req)
	go SendKafkaData(req)
	if err != nil {
		return err
	}
	//fmt.Println(string(v))
	rsp.Msg = string(v)
	return nil
}

// SendKafkaData 发送消息到kafka中
func SendKafkaData(req *uavdata.Request) {
	vv,err := json.Marshal(req)
	if err != nil {
		fmt.Println("json Marshal, err:", err)
		return
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll          // 发送完数据需要leader和follow都确认
	config.Producer.Partitioner = sarama.NewRandomPartitioner // 新选出一个partition
	config.Producer.Return.Successes = true                   // 成功交付的消息将在success channel返回

	// 构造一个消息
	msg := &sarama.ProducerMessage{}
	msg.Topic = "web_log"
	msg.Value = sarama.StringEncoder(vv)
	// 连接kafka
	client, err := sarama.NewSyncProducer([]string{"127.0.0.1:9092"}, config)
	if err != nil {
		fmt.Println("producer closed, err:", err)
		return
	}
	defer client.Close()
	// 发送消息
	pid, offset, err := client.SendMessage(msg)
	if err != nil {
		fmt.Println("send msg failed, err:", err)
		return
	}
	fmt.Printf("pid:%v offset:%v\n", pid, offset)


}