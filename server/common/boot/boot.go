package boot

import (
	"flag"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"hotgo/common/apollo"
)

var ctx = gctx.GetInitCtx()

func init() {
	configFile := flag.String("f", "./manifest/config/config.yaml", "The config file path")
	flag.Parse()
	// 设置默认配置文件
	g.Cfg().GetAdapter().(*gcfg.AdapterFile).SetFileName(*configFile)

	adapter, err := apollo.New(ctx, apollo.Config{
		AppID:             g.Cfg().MustGet(ctx, `apollo.AppID`).String(),
		IP:                g.Cfg().MustGet(ctx, `apollo.Ip`).String(),
		Cluster:           g.Cfg().MustGet(ctx, `apollo.Cluster`).String(),
		NamespaceName:     g.Cfg().MustGet(ctx, `apollo.NamespaceName`).String(),
		IsBackupConfig:    g.Cfg().MustGet(ctx, `apollo.IsBackupConfig`).Bool(),
		BackupConfigPath:  g.Cfg().MustGet(ctx, `apollo.BackupConfigPath`).String(),
		Secret:            g.Cfg().MustGet(ctx, `apollo.Secret`).String(),
		SyncServerTimeout: g.Cfg().MustGet(ctx, `apollo.SyncServerTimeout`).Int(),
		MustStart:         g.Cfg().MustGet(ctx, `apollo.MustStart`).Bool(),
		Watch:             g.Cfg().MustGet(ctx, `apollo.Watch`).Bool(),
	})
	if err != nil {
		g.Log().Fatalf(ctx, `%+v`, err)
	}
	//将apollo适配器添加到默认配置中
	g.Cfg().SetAdapter(adapter)

	//获取MySQL账号密码
	configNode := gdb.ConfigNode{
		Host:   "127.0.0.1",
		Port:   "3306",
		User:   "hotgo",
		Pass:   "hotgo123456.",
		Link:   "mysql:hotgo:hg123456.@tcp(127.0.0.1:3306)/hotgo?loc=Local&parseTime=true&charset=utf8mb4",
		Prefix: "hg_",
	}
	//手动设置新的mysql配置
	gdb.AddConfigNode("default", configNode)

	link := "mysql:hotgo:hg123456.@tcp(127.0.0.1:3306)/hotgo?loc=Local&parseTime=true&charset=utf8mb4"
	//动态修改配置内容
	dynamicContent := map[string]interface{}{
		"database.default.link":   link,
		"database.default.Prefix": "hg_",
	}
	adapter.(*apollo.Client).SetContent(dynamicContent)
}
