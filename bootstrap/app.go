package bootstrap

import "herostory-server/pkg/logger"

func InitApp() {
	logger.InitZeroLogger("./storage/logs", "biz_server")
}
