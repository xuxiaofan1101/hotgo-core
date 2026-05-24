package genrouter

import "hotgo/internal/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.DataFieldTemplate) // 字段模板
}
