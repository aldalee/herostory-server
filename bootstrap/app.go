package bootstrap

import (
	"herostory-server/internal/pb"
	"herostory-server/pkg/logger"
)

func InitApp() {
	logger.InitZeroLogger("./storage/logs", "biz_server")
	pb.InitMaps()
}
