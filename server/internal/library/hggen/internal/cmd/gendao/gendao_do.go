// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gendao

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gview"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"

	"hotgo/internal/library/hggen/internal/consts"
	"hotgo/internal/library/hggen/internal/utility/mlog"
	"hotgo/internal/library/hggen/internal/utility/utils"
)

func generateDo(ctx context.Context, in CGenDaoInternalInput) {
	var dirPathDo = filepath.FromSlash(gfile.Join(in.Path, in.DoPath))
	in.genItems.AppendDirPath(dirPathDo)
	in.NoJsonTag = true
	in.DescriptionTag = false
	in.NoModelComment = false
	// Model content.
	for i, tableName := range in.TableNames {
		fieldMap, err := in.DB.TableFields(ctx, tableName)
		if err != nil {
			mlog.Fatalf("fetching tables fields failed for table '%s':\n%v", tableName, err)
		}
		var (
			newTableName        = in.NewTableNames[i]
			doFilePath          = gfile.Join(dirPathDo, gstr.CaseSnake(newTableName)+".go")
			structDefinition, _ = generateStructDefinition(ctx, generateStructDefinitionInput{
				CGenDaoInternalInput: in,
				TableName:            tableName,
				StructName:           formatFieldName(newTableName, FieldNameCaseCamel),
				FieldMap:             fieldMap,
				IsDo:                 true,
			})
		)
		// replace all types to interface{}.
		structDefinition, _ = gregex.ReplaceStringFuncMatch(
			"([A-Z]\\w*?)\\s+([\\w\\*\\.]+?)\\s+(//)",
			structDefinition,
			func(match []string) string {
				// If the type is already a pointer/slice/map, it does nothing.
				if !gstr.HasPrefix(match[2], "*") && !gstr.HasPrefix(match[2], "[]") && !gstr.HasPrefix(match[2], "map") {
					return fmt.Sprintf(`%s interface{} %s`, match[1], match[3])
				}
				return match[0]
			},
		)
		modelContent := generateDoContent(
			ctx,
			in,
			tableName,
			formatFieldName(newTableName, FieldNameCaseCamel),
			structDefinition,
		)
		in.genItems.AppendGeneratedFilePath(doFilePath)
		err = gfile.PutContents(doFilePath, strings.TrimSpace(modelContent))
		if err != nil {
			mlog.Fatalf(`writing content to "%s" failed: %v`, doFilePath, err)
		} else {
			utils.GoFmt(doFilePath)
			mlog.Print("generated:", gfile.RealPath(doFilePath))
		}
	}
}

func generateDoContent(
	ctx context.Context, in CGenDaoInternalInput, tableName, tableNameCamelCase, structDefine string,
) string {
	var (
		tplContent = getTemplateFromPathOrDefault(
			in.TplDaoDoPath, consts.TemplateGenDaoDoContent,
		)
	)
	tplView.ClearAssigns()
	tplView.Assigns(gview.Params{
		tplVarTableName:          tableName,
		tplVarPackageImports:     getImportPartContent(ctx, structDefine, true, nil),
		tplVarTableNameCamelCase: tableNameCamelCase,
		tplVarStructDefine:       structDefine,
		tplVarPackageName:        filepath.Base(in.DoPath),
	})
	assignDefaultVar(tplView, in)
	doContent, err := tplView.ParseContent(ctx, tplContent)
	if err != nil {
		mlog.Fatalf("parsing template content failed: %v", err)
	}
	return doContent
}
