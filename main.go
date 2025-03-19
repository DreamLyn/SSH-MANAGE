package main

import (
	"flag"
	"log"
	"log/slog"
	"os"
	"strings"

	"github.com/dreamlyn/ssh-manage/internal/app"
	"github.com/dreamlyn/ssh-manage/internal/rest/handlers"

	"github.com/dreamlyn/ssh-manage/ui"
	"github.com/pocketbase/pocketbase/apis"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/plugins/migratecmd"
	"github.com/pocketbase/pocketbase/tools/hook"

	// "github.com/pocketbase/pocketbase/ui"

	_ "github.com/dreamlyn/ssh-manage/migrations"
)

func main() {
	app := app.GetApp()

	// 获取运行参数
	var flagHttp string
	var flagDir string
	flag.StringVar(&flagHttp, "http", "127.0.0.1:8090", "HTTP server address")
	flag.StringVar(&flagDir, "dir", "/pb_data/database", "Pocketbase data directory")
	if len(os.Args) < 2 {
		slog.Error("[ssh-manage] missing exec args")
		os.Exit(1)
		return
	}
	_ = flag.CommandLine.Parse(os.Args[2:]) // skip the first two arguments: "main.go serve"

	// loosely check if it was executed using "go run"
	isGoRun := strings.HasPrefix(os.Args[0], os.TempDir())
	migratecmd.MustRegister(app, app.RootCmd, migratecmd.Config{
		// enable auto creation of migration files when making collection changes in the Dashboard
		// (the isGoRun check is to enable it only during development)
		Automigrate: isGoRun,
	})

	//监听系统启动，初始化
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		return e.Next()
	})

	//websocket终端通信
	sshHandler := handlers.NewSSHHandler()
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.GET("/ws/{id}", func(e *core.RequestEvent) error {
			sshHandler.Handle(e.Response, e.Request) // 将 e.Response 和 e.Request 传递给 handler.Handle
			return nil
		})
		return e.Next()
	})
	//端口转发支持
	forwardingHandler := handlers.NewPortForwardingHandler()
	app.OnServe().BindFunc(func(e *core.ServeEvent) error {
		e.Router.POST("/forwarding/{id}", func(e *core.RequestEvent) error {
			forwardingHandler.Handle(e.Response, e.Request) // 将 e.Response 和 e.Request 传递给 handler.Handle
			return nil
		})
		return e.Next()
	})

	app.OnServe().Bind(&hook.Handler[*core.ServeEvent]{
		Func: func(e *core.ServeEvent) error {
			e.Router.
				GET("/{path...}", apis.Static(ui.DistDirFS, false)).
				Bind(apis.Gzip())
			return e.Next()
		},
		Priority: 999,
	})

	app.OnTerminate().BindFunc(func(e *core.TerminateEvent) error {
		// routes.Unregister()
		slog.Info("[SSH-MANAGE] Exit!")
		return e.Next()
	})

	slog.Info("[SSH-MANAGE] Visit the website: http://" + flagHttp)

	// start pocketbase
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
