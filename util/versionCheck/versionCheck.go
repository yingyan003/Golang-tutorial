package versionCheck

import (
	"strings"
	"strconv"
	"fmt"
)

type VersionInfo struct {
	//所有版本，为1表示作用于所有版本，为0则以下面指定的版本范围为准
	AllVersion int8 `json:"allVersion"`
	//最小版本
	MinVersion string `json:"minVersion"`
	//是否包含最小版本，即是开区间还是闭区间
	MinIncluded int8 `json:"minIncluded"`
	//最大版本
	MaxVersion string `json:"maxVersion"`
	//是否包含最大版本
	MaxIncluded int8 `json:"maxIncluded"`
}

//检查2个配置版本是否交叉
//true表示无交叉，反之亦然
func CheckVersionConflict(v1, v2 *VersionInfo) bool {
	//参数校验
	if v1 == nil || v2 == nil {
		fmt.Errorf("版本检查参数错误：VersionInfo不能为空")
		return false
	}
	//版本为空默认为正/负无穷。
	//在allVersion==0时（即不适用于所有版本），一个配置的最大最小版本不能同时空(为空表示版本范围为正负无穷）
	if ( 0 == v1.AllVersion && (v1.MinVersion == "" && v1.MaxVersion == "")) || (0 == v2.AllVersion && (v2.MinVersion == "" && v2.MaxVersion == "")) {
		fmt.Errorf("版本检查参数错误，最小最大版本号不能同时为空")
		return false
	}
	//若有一个配置作用于所有版本，有交叉
	if v1.AllVersion == 1 || v2.AllVersion == 1 {
		return false
	}
	//正负无穷检查
	//检查一个配置版本范围有一边是无穷的情况
	minInfinityCheck(v1, v2)
	maxInfinityCheck(v1, v2)

	//将string类型的版本转换为int类型。ints[0-3]分别对应v1min,v1max,v2min,v2max
	ints := toInt(v1.MinVersion, v1.MaxVersion, v2.MinVersion, v2.MaxVersion)

	//参数校验
	if (ints[0] == ints[1] && (v1.MinIncluded == 0 || v1.MaxIncluded == 0)) || (ints[2] == ints[3] && (v2.MinIncluded == 0 || v2.MaxIncluded == 0)) {
		fmt.Errorf("版本检查参数错误：最小最大版本相同时，边界必须都包含")
		return false
	}
	//版本边界检查
	ints[0] = minIncludeCheck(ints[0], v1.MinIncluded)
	ints[1] = maxIncludeCheck(ints[1], v1.MaxIncluded)
	ints[2] = minIncludeCheck(ints[2], v2.MinIncluded)
	ints[3] = maxIncludeCheck(ints[3], v2.MaxIncluded)

	//无交叉
	if ints[1] < ints[2] || ints[3] < ints[0] {
		return true
	}
	//恢复正负无穷的初始值""
	resetInfinity(v1, v2)
	return false
}

//将负无穷置为"0"
func minInfinityCheck(minV ...*VersionInfo) {
	for _, v := range minV {
		if v.MinVersion == "" {
			v.MinVersion = "0"
		}
	}
}

//将正无穷置为"1000000"
func maxInfinityCheck(maxV ...*VersionInfo) {
	for _, v := range maxV {
		if v.MaxVersion == "" {
			v.MaxVersion = "1000000"
		}
	}
}

//版本格式转化：1.每个部分不足2位用0补足 2.string转为int
//按参数传入的顺序依次返回转化后的数据
func toInt(strs ...string) (ints []int) {
	for _, str := range strs {
		strSum := ""
		ss := strings.Split(str, ".")
		for _, s := range ss {
			//fmt.Println("i,len()",i,	len(s))
			if len(s) == 1 {
				strSum = strSum + "0" + s
			} else {
				strSum = strSum + s
			}
		}
		i, _ := strconv.Atoi(strSum)
		ints = append(ints, i)
	}
	return ints
}

func minIncludeCheck(minV int, incl int8) int {
	//当minV为负无穷时不操作
	if minV != 0 && incl == 0 {
		return minV + 1
	}
	return minV
}

func maxIncludeCheck(maxV int, incl int8) int {
	//当maxV是正无穷不操作
	if maxV != 1000000 && incl == 0 {
		return maxV - 1
	}
	return maxV
}

//前面将负无穷置为0，正无穷置为1000000。本方法恢复其初始值""
func resetInfinity(vs ...*VersionInfo) {
	for _, v := range vs {
		if v.MinVersion == "0" {
			v.MinVersion = ""
		}
		if v.MaxVersion == "1000000" {
			v.MinVersion = ""
		}
	}
}
