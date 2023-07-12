package crawling

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/go-redis/redis/v8"
)

var comandHandlers = make(map[string]CommandHandler)

type (
	CommandSchema struct {
		Command string   `json:"command"`
		Args    []string `json:"args"`
	}

	CommandResponse struct {
		CommandSchema
		Response string `json:"response"`
	}

	CommandHandler func(command CommandSchema) string
)

func runCommandHandler(ctx context.Context) {
	pb := cli.Subscribe(ctx, "crawler_command")
	channel := pb.Channel()
	logger.Infof("正在啟動 redis 指令監聽...")
	defer func() {
		if err := pb.Close(); err != nil {
			logger.Warnf("停止指令監聽時出現錯誤: %v", err)
		}
		logger.Infof("指令監聽已停止。")
	}()
	defer func() {
		if err := recover(); err != nil {
			logger.Errorf("處理指令時出現致命錯誤: %v", err)
			debug.PrintStack()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			logger.Debugf("收到中止指令，正在停止指令監聽")
			return

		case cmd, ok := <-channel:
			if !ok {
				logger.Debugf("訂閱接收閘口關閉，正在停止指令監聽")
				return
			}
			go handleCommand(cmd)
		}
	}
}

func handleCommand(msg *redis.Message) {
	var command CommandSchema
	err := json.Unmarshal([]byte(msg.Payload), &command)
	if err != nil {
		logger.Errorf("解析指令時出現錯誤: %v", err)
		return
	}
	logger.Infof("收到指令: %s (%s)", command.Command, strings.Join(command.Args, ", "))

	var res string

	if handler, ok := comandHandlers[command.Command]; ok {
		res = handler(command)
	} else {
		res = fmt.Sprintf("未知的指令: %s", command.Command)
	}

	if res == "" {
		return
	}

	logger.Infof("指令回應: %s", res)
	cli.Publish(ctx, "crawler_command_response", &CommandResponse{
		CommandSchema: command,
		Response:      res,
	})
}

func AddCommandHandler(command string, handler CommandHandler) {
	comandHandlers[command] = handler
}
