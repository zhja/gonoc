package gonoc

import (
    //"fmt"
    "strings"
)

//string
//字符串截取
func Substr(str string, start, length int) string {
    rs := []rune(str)
    rl := len(rs)
    end := 0

    if start < 0 {
        start = rl - 1 + start
    }
    end = start + length

    if start > end {
        start, end = end, start
    }

    if start < 0 {
        start = 0
    }
    if start > rl {
        start = rl
    }
    if end < 0 {
        end = 0
    }
    if end > rl {
        end = rl
    }

    return string(rs[start:end])
}

//首字母大写 小写
func StringFC(str string) (s string) {
    strSlice := strings.Split(str, "")
    for key, val := range strSlice {
        if key == 0 {
            strSlice[key] = strings.ToUpper(val)
        }
    }
    s = strings.Join(strSlice, "")
    return
}

//判断slice中值是否存在(一维)
func ExistsSV(slices []string, value string) bool {
    var t = false
    for _, val := range slices {
        if val == value {
            t = true
            break
        }
    }
    return t
}

//判断slice中值是否存在(二维)
func ExistsSVT(slices [][]string, value string) (t bool, keys int) {
    for key, val := range slices {
        if val[0] == value {
            t = true
            keys = key
            break
        }
    }
    return
}