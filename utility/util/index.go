package util

import (
	"context"
	"github.com/bwmarrin/snowflake"
	"github.com/gogf/gf/v2/frame/g"
)

// 雪花
func CreateSnowflake(ctx context.Context) *snowflake.Node {
	node, err := snowflake.NewNode(1) // 1 是节点的ID
	if err != nil {
		g.Log().Error(ctx, err.Error())
	}

	return node
}
