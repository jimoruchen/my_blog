package snowflake

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

var node *snowflake.Node

// Init 初始化 Snowflake 节点
func Init(startTime string, machineID int64) (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return err
	}
	// 设置自定义纪元（epoch），单位为毫秒
	snowflake.Epoch = st.UnixNano() / 1e6 // 转换为毫秒
	// 创建节点
	node, err = snowflake.NewNode(machineID)
	return err
}

// GenID 生成唯一 ID
func GenID() int64 {
	return node.Generate().Int64()
}

//func main() {
//	// 初始化：以 "2020-07-01" 为起始时间，机器 ID 为 1
//	if err := Init("2020-07-01", 1); err != nil {
//		fmt.Printf("init failed, err:%v\n", err)
//		return
//	}
//	// 生成一个 ID
//	id := GenID()
//	fmt.Println(id)
//}
