package main

import (
	"flag"
	"fmt"
	"github.com/jixindatech/sqlaudit/audit/server"
	"github.com/jixindatech/sqlaudit/pkg/config"
	"github.com/jixindatech/sqlaudit/pkg/core/alert"
	"github.com/jixindatech/sqlaudit/pkg/core/golog"
	"github.com/jixindatech/sqlaudit/pkg/core/monitor"
	"github.com/jixindatech/sqlaudit/pkg/queue"
	"github.com/jixindatech/sqlaudit/pkg/storage"
	"github.com/jixindatech/sqlaudit/pkg/task"
	"github.com/jixindatech/sqlaudit/pkg/webserver"
	"github.com/jixindatech/sqlaudit/pkg/webserver/models"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var configFile *string = flag.String("config", "./etc/config.yaml", "kingshard config file")
var capture *bool = flag.Bool("c", false, "the capture way")
var inf *string = flag.String("i", "eth0", "capture interface, default eth0")
var version *bool = flag.Bool("v", false, "the version ")

var (
	BuildDate    string
	BuildVersion string
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Parse()

	if *version {
		fmt.Printf("Git commit:%s\n", BuildVersion)
		fmt.Printf("Build time:%s\n", BuildDate)
		return
	}

	_ = golog.SetDefaultZapLog()

	if *capture && len(*inf) == 0 {
		golog.Fatal("main", zap.String("err", "you must specify an inerface to capture"))
	}

	if len(*configFile) == 0 {
		golog.Fatal("main", zap.String("err", "you must use a config file"))
	}

	var err error
	cfg, err := config.ParseConfigFile(*configFile)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	if len(cfg.LogPath) != 0 {
		_ = golog.InitZapLog(cfg.LogLevel, cfg.LogPath)
	} else {
		_ = golog.InitZapLog(cfg.LogLevel, "stdout")
	}

	//Queue init phase
	queueInstance, err := queue.GetQueue(cfg.QueueType, cfg.QueueNum)
	if queueInstance == nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	//Storage init phase
	storageInstance, err := storage.GetStorage(cfg.EsConfig)
	if storageInstance == nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}
	//Model opened
	err = models.OpenDatabase(cfg.Database, storageInstance)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	apiSvr, err := webserver.NewApiServer(cfg)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	svr, err := server.NewServer(cfg, *capture, *inf, queueInstance)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	emailAlert, err := alert.GetAlert(cfg)
	if emailAlert == nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}

	var mysqlTask task.Task
	mysqlTask, err = task.GetTask("mysql", "mysql", queueInstance, emailAlert)
	if mysqlTask == nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}
	go mysqlTask.Run()

	prom, err := monitor.NewPrometheus("0.0.0.0:9001", svr)
	if err != nil {
		golog.Fatal("main", zap.String("err", err.Error()))
	}
	go prom.Run()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGPIPE,
		//syscall.SIGUSR1,
	)

	go apiSvr.Run()
	go svr.Run()

	golog.Info("main", zap.String("status", "running"))

	for {
		sig := <-sc
		if sig == syscall.SIGINT || sig == syscall.SIGTERM || sig == syscall.SIGQUIT {
			golog.Info("main", zap.String("signal", fmt.Sprintf("%d", sig)))
			svr.Close()
			golog.Close()
			break
		} else if sig == syscall.SIGPIPE {
			//golog.Info("main", "main", "Ignore broken pipe signal", 0)
			// skip
		}
	}

	golog.Info("main", zap.String("status", "quit"))
}
