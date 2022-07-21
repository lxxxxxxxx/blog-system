package main

import (
	"context"
	"flag"
	"fmt"

	"lxxxxxxxx.github.com/applet/backend/pkg/api"
	"lxxxxxxxx.github.com/applet/backend/pkg/common"
	"lxxxxxxxx.github.com/applet/backend/pkg/models"
)

var (
	configPath = flag.String("c", "./etc/benckend.yml", "config file path.")
)

func main() {
	flag.Parse()
	fmt.Println("config path:", *configPath)

	if err := common.InitConfig(*configPath); err != nil {
		fmt.Errorf("init config failed,error:%w", err)
	}
	ctx := context.Background()

	models.Setup(ctx)

	s := api.NewServer()
	s.RunGrpcAndGateway(ctx)
}
