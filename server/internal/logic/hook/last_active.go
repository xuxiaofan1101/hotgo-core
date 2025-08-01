// Package hook
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package hook

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/utility/simple"
	"sync"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
)

type visitor struct {
	lastSeen *gtime.Time
}

var (
	visitors = make(map[int64]*visitor)
	mtx      sync.Mutex
)

// Run a background goroutine to remove old entries from the visitors map.
func init() {
	go cleanupVisitors()
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func cleanupVisitors() {
	gtimer.AddSingleton(gctx.New(), time.Minute, func(ctx context.Context) {
		mtx.Lock()
		for memberId, v := range visitors {
			if gtime.Now().Sub(v.lastSeen) > 3*time.Minute {
				delete(visitors, memberId)
			}
		}
		mtx.Unlock()
	})
}

func allow(memberId int64) bool {
	mtx.Lock()
	defer mtx.Unlock()
	v, exists := visitors[memberId]
	if !exists {
		visitors[memberId] = &visitor{gtime.Now()}
		return true
	}

	if gtime.Now().Sub(v.lastSeen) > time.Second*3 {
		v.lastSeen = gtime.Now()
		return true
	}
	return false
}

// lastAdminActive 更新后台用户最后活跃
func (s *sHook) lastAdminActive(r *ghttp.Request) {
	if r.IsFileRequest() {
		return
	}

	var (
		ctx    = r.Context()
		member = contexts.GetUser(ctx)
	)

	if member == nil || member.App != consts.AppAdmin {
		return
	}

	if allow(member.Id) {
		simple.SafeGo(ctx, func(ctx context.Context) {
			_, err := dao.AdminMember.Ctx(ctx).
				Where(dao.AdminMember.Columns().Id, member.Id).
				Data(g.Map{dao.AdminMember.Columns().LastActiveAt: gtime.Now()}).
				Update()
			if err != nil {
				g.Log().Warningf(ctx, "hook lastActive err:%+v, memberId:%v", err, member.Id)
			}
		})
	}
}
