package service

import "strconv"

// itoa uint64 → string，用于 JWT sub 与缓存 key。
func itoa(v uint64) string { return strconv.FormatUint(v, 10) }
