package initialize

import (
	"fmt"
	"github.com/DouYin/service/global"
	"github.com/bwmarrin/snowflake"
)

func Snowflake() {
	node, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	global.ID = node
}
