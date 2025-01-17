package service

import (
	"demo1/src/model"
	"github.com/spf13/cast"
	"math/rand"
	"strings"
)

func Hit(form model.Client) (string, string, string, string, string) {
	//todo test
	//if hit,
	//return downloadUrl, updateVersionCode, md5, title, updateTips

	var downloadUrl, updateVersionCode, md5, title, updateTips string

	//需要 model 层实现返回规则集的接口 已完成
	//var rules *[]model.Rule
	rules := model.GetRules(form.DeviceId)
	//order by updateVersionCode, use quickSort O(nlogn)
	quickSort(rules, 0, len(*rules)-1)

	//O(n)
	for _, rule := range *rules {
		if matchRule(&rule, &form) {
			downloadUrl, updateVersionCode, md5, title, updateTips = getDownloadInfo(&rule)
			break
		}
	}
	return downloadUrl, updateVersionCode, md5, title, updateTips
}

func quickSort(rules *[]model.Rule, l, r int) {
	if l >= r {
		return
	}

	mid := l
	randIdx := rand.Intn(r-l+1) + l
	(*rules)[l], (*rules)[randIdx] = (*rules)[randIdx], (*rules)[l]
	target := (*rules)[l].UpdateVersionCode
	for i, rule := range (*rules)[l : r+1] {
		if rule.UpdateVersionCode <= target {
			(*rules)[mid], (*rules)[l+i] = (*rules)[l+i], (*rules)[mid]
			mid++
		}
	}
	//左侧 小于等于 target
	//右侧 大于 target
	quickSort(rules, l, mid-1)
	quickSort(rules, mid, r)
	return
}

func matchRule(rule *model.Rule, form *model.Client) bool {
	//model.Client.Version : 请求api版本，⽐如v1/v2
	//model.Client.version_code : 应⽤⼤版本，⽐如8.1.4
	//deviceIdList 白名单，model 里处理
	if form.DevicePlatform != rule.Platform {
		//设备平台
		return false
	}
	if form.Aid != rule.Aid {
		//app 是否相同
		return false
	}
	if (*rule).Channel != (*form).Channel {
		//渠道 是否相同
		return false
	}

	//是否符合 版本要求（应⽤⼩版本，⽐如8.1.4.01）
	maxVersionCode := strings.Split((*rule).MaxUpdateVersionCode, ".")
	minVersionCode := strings.Split((*rule).MinUpdateVersionCode, ".")
	versionCode := strings.Split((*form).UpdateVersionCode, ".")
	for index := 0; index < 4; index++ {
		if cast.ToInt(versionCode[index]) < cast.ToInt(minVersionCode[index]) || cast.ToInt(versionCode[index]) > cast.ToInt(maxVersionCode[index]) {
			return false
		}
		if cast.ToInt(versionCode[index]) != cast.ToInt(minVersionCode[index]) || cast.ToInt(versionCode[index]) != cast.ToInt(maxVersionCode[index]) {
			break
		}
	}

	if (*form).OsApi < (*rule).MinOsApi || (*form).OsApi > (*rule).MaxOsApi {
		//系统 是否适配
		return false
	}
	if (*form).CpuArch != (*rule).CpuArch {
		return false
	}
	return true
}

func getDownloadInfo(rule *model.Rule) (string, string, string, string, string) {
	return (*rule).DownloadUrl, (*rule).UpdateTips, (*rule).Md5, (*rule).Title, (*rule).UpdateTips
}
