package data

import (
	"go-micro.dev/v5/logger"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-micro-example/app/user/service/internal/data/models"

	"go-micro-example/api/gen/go/common/conf"
)

// NewGormClient 创建数据库客户端
func NewGormClient(cfg *conf.Data, l logger.Logger) *gorm.DB {
	if cfg == nil {
		return nil
	}

	var driver gorm.Dialector
	switch cfg.Database.Driver {
	default:
		fallthrough
	case "mysql":
		driver = mysql.Open(cfg.Database.Source)
		break
	case "postgres":
		driver = postgres.Open(cfg.Database.Source)
		break
		//case "clickhouse":
		//	driver = clickhouse.Open(cfg.Database.Source)
		//	break
		//case "sqlite":
		//	driver = sqlite.Open(cfg.Database.Source)
		//	break
		//case "sqlserver":
		//	driver = sqlserver.Open(cfg.Database.Source)
		//break
	}

	client, err := gorm.Open(driver, &gorm.Config{})
	if err != nil {
		l.Logf(logger.FatalLevel, "failed opening connection to db: %v", err)
		return nil
	}

	// 运行数据库迁移工具
	if cfg.Database.Migrate {
		if err := client.AutoMigrate(
			models.GetMigrates()...,
		); err != nil {
			l.Logf(logger.FatalLevel, "failed creating schema resources: %v", err)
			return nil
		}
	}

	return client
}
