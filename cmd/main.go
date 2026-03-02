package main

import (
	"billohub/config"
	"billohub/internal/api"
	"billohub/internal/manager"
	"billohub/internal/model"
	"billohub/internal/skill"
	"billohub/internal/storage"
	gateway "billohub/pkg/gateway"
	"billohub/pkg/helper"
	"billohub/pkg/logx"
	middleware "billohub/pkg/middleware"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const pidFile = "billo-hub.pid"

func main() {
	// --- 1. PID   ---
	if err := writePIDFile(pidFile); err != nil {
		fmt.Printf("Failed to write PID file: %v\n", err)
		os.Exit(1)
	}
	defer func() {
		if err := removePIDFile(pidFile); err != nil {
			fmt.Printf("Failed to remove PID file: %v\n", err)
		}
	}()

	// --- 2. ctx---
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		fmt.Printf("Received signal: %s. Shutting down...\n", sig)
		cancel()
	}()

	// --- 3. Init ---
	config.Init()
	cfg := config.GetConfig()

	logger, _, err := logx.New(cfg.LogConf)
	if err != nil {
		fmt.Println("init log error", err.Error())
		return
	}

	httpConf := cfg.LogConf
	httpConf.JsonLog.Path = httpConf.JsonLog.Path + "_http"
	httpConf.CsvLog.Path = httpConf.CsvLog.Path + "_http"
	httpLogger, _, err := logx.New(httpConf)
	if err != nil {
		logger.Error("init http log error", zap.Error(err))
		return
	}
	var storagePtr model.AgentStorage

	storagePtr, err = storage.NewStorage(cfg.DSN)
	if err != nil {
		logger.Error("init storage error", zap.Error(err))
		return
	}

	// --- 数据库和默认用户初始化 ---
	if err := initDatabase(storagePtr.DB(), logger); err != nil {
		logger.Error("failed to initialize database", zap.Error(err))
		return
	}
	if err := initDefaultUser(storagePtr, logger); err != nil {
		logger.Error("failed to initialize default user", zap.Error(err))
		return
	}

	globalSkills := skill.GlobalSKills{}

	hub := manager.NewAgentHub(storagePtr, globalSkills, logger)
	err = hub.InitLoad()
	if err != nil {
		logger.Error("hub init load err", zap.Error(err))
		return
	}

	// --- 4. 启动 Gin Web 服务器 ---
	ginEngine := gateway.NewEngine(cfg.DebugMode)
	ginEngine.Use(middleware.RecoveryWithZap(false, httpLogger), middleware.HttpLog(false, httpLogger))

	// 注册核心 API 路由
	api.RegisterRoutes(ginEngine, cfg.DebugMode, hub, logger)
	ginEngine.Static("/web", "./web")
	ginEngine.NoRoute(func(c *gin.Context) {
		// 检查请求路径是否以 /web/ 开头
		if c.Request.URL.Path[:5] == "/web/" {
			// 返回 index.html 文件
			c.File("./web/index.html")
			return
		}
		// 对于其他所有未匹配的路由，可以返回标准的404页面
		c.JSON(http.StatusNotFound, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})
	logger.Info("Starting, listen on " + cfg.Http)

	gw := gateway.NewGateway(cfg.Http, ginEngine, logger)

	if err := gw.Start(ctx); err != nil {
		logger.Error("HTTP server stopped with error", zap.Error(err))
	} else {
		logger.Info("HTTP server exited gracefully.")
	}
}

// initDatabase uses GORM's AutoMigrate to initialize the database schema.
func initDatabase(db *gorm.DB, logger *zap.Logger) error {
	logger.Info("Auto-migrating database schema...")
	// GORM will automatically create tables, missing columns, and indexes.
	// It will not delete unused columns to protect your data.
	err := db.AutoMigrate(
		&model.User{},
		&model.AgentInstanceData{},
		&model.Message{},
		&model.Chat{},
		&model.ScheduledTask{},
		&model.LLMModel{},
	)
	if err != nil {
		return helper.WrapError(err, "gorm auto-migrate failed")
	}

	logger.Info("Database schema migration completed successfully.")
	return nil
}

// initDefaultUser checks if the default admin user exists, and creates it if not.
func initDefaultUser(storageSvc model.AgentStorage, logger *zap.Logger) error {
	_, err := storageSvc.GetUserByUsername("admin")
	if err == nil {
		logger.Info("Default admin user already exists.")
		return nil // User exists, do nothing
	}

	logger.Info("Default admin user not found, creating one...")
	hashedPassword, err := helper.HashPassword("123456")
	if err != nil {
		return helper.WrapError(err, "failed to hash default password")
	}

	defaultUser := &model.User{
		Username:       "admin",
		HashedPassword: hashedPassword,
	}

	if err := storageSvc.CreateUser(defaultUser); err != nil {
		return helper.WrapError(err, "failed to create default admin user")
	}
	logger.Info("Default admin user created successfully with password '123456'. Please change it immediately.")

	return nil
}

func writePIDFile(filename string) error {
	pid := os.Getpid()
	return os.WriteFile(filename, []byte(fmt.Sprintf("%d", pid)), 0644)
}

func removePIDFile(filename string) error {
	return os.Remove(filename)
}
