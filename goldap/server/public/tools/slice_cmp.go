package tools

import (
	"fmt"
	"strconv"
	"strings"
)

// ArrStrCmp compares two string slices and returns added/deleted elements
func ArrStrCmp(src []string, dest []string) ([]string, []string) {
	msrc := make(map[string]byte)
	mall := make(map[string]byte)
	var set []string
	
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) {
			continue
		}
		set = append(set, v)
	}
	for _, v := range set {
		delete(mall, v)
	}
	
	var added, deleted []string
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}

// ArrUintCmp compares two uint slices and returns added/deleted elements
func ArrUintCmp(src []uint, dest []uint) ([]uint, []uint) {
	msrc := make(map[uint]byte)
	mall := make(map[uint]byte)
	var set []uint
	
	for _, v := range src {
		msrc[v] = 0
		mall[v] = 0
	}
	for _, v := range dest {
		l := len(mall)
		mall[v] = 1
		if l != len(mall) {
			continue
		}
		set = append(set, v)
	}
	for _, v := range set {
		delete(mall, v)
	}
	
	var added, deleted []uint
	for v := range mall {
		_, exist := msrc[v]
		if exist {
			deleted = append(deleted, v)
		} else {
			added = append(added, v)
		}
	}
	return added, deleted
}

// SliceToString converts uint slice to delimited string
func SliceToString(src []uint, delim string) string {
	return strings.Trim(strings.ReplaceAll(fmt.Sprint(src), " ", delim), "[]")
}

// StringToSlice converts delimited string to uint slice
func StringToSlice(src string, delim string) []uint {
	var dest []uint
	if src == "" {
		return dest
	}
	strs := strings.Split(src, delim)
	for _, v := range strs {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		t, err := strconv.Atoi(v)
		if err == nil {
			dest = append(dest, uint(t))
		}
	}
	return dest
}
