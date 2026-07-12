package handler

import (
	"strconv"
	"strings"

	"github.com/TDroyal/Accounting/server/internal/service"
)

// categoryNameMap 把分类树扁平化为 {id: name} 映射，供列表/导出补全分类名。
func categoryNameMap(tree []service.CategoryNode) map[uint64]string {
	m := map[uint64]string{}
	for _, n := range tree {
		m[n.ID] = n.Name
		for _, ch := range n.Children {
			m[ch.ID] = n.Name + "/" + ch.Name
		}
	}
	return m
}

// toCSV 将导出行拼为 CSV 字符串（含表头）。
func toCSV(rows []service.ExportRow) string {
	var b strings.Builder
	b.WriteString("occurred_at,type,category,amount,note\n")
	for _, r := range rows {
		b.WriteString(csvField(r.OccurredAt) + ",")
		b.WriteString(csvField(typeName(r.Type)) + ",")
		b.WriteString(csvField(r.CategoryName) + ",")
		b.WriteString(csvField(strconv.FormatFloat(r.Amount, 'f', 2, 64)) + ",")
		b.WriteString(csvField(r.Note) + "\n")
	}
	return b.String()
}

// csvField 简单 CSV 字段转义：含逗号/引号/换行则加引号并转义内部引号。
func csvField(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

func typeName(t int8) string {
	if t == 1 {
		return "转账"
	}
	return "支出"
}

