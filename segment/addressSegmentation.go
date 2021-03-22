package segment

import (
	"bytes"
	"canOnline/poiUtil"
	//"fmt"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

// AddressSegmentation 中文地址分割
type AddrSegmentation struct {
	poiUtil.XingZhengQuHuaDict
	index            int
	pXzqh            *poiUtil.XingZhengQuHuaItem
	hit              int32
	filling          byte // 原地址要素识别以后的填充符
	xzqhSep          rune // 行政区划字典表 mergename 中的地址要素分隔符
	extraInfoSep     byte // 剩余地址信息部分的分割符 ‘,’
	roadRegex        *regexp.Regexp
	roadNumRegex     *regexp.Regexp
	buildingNumRegex *regexp.Regexp
	poiRegex         *regexp.Regexp
	unitNumRegex     *regexp.Regexp
	floorNumRegex    *regexp.Regexp
	roomNumRegex     *regexp.Regexp
}

// isInRegion 判断当前区域是否在上个区域所辖范围内
// 区县以上 地址编码是两位表示一个区域 13 06 38 103215
// 区县以下 地址编码是三位表示一个区域 130638 103 215
func (as *AddrSegmentation) isInRegion(last *poiUtil.XingZhengQuHuaItem, paddrCode *string, currentLevel uint8) bool {
	inRegion := false
	if currentLevel-last.Level >= 2 {
		matchLen := 0
		switch last.Level {
		case poiUtil.Province, poiUtil.City, poiUtil.District:
			matchLen = int(last.Level*2 + 2)
		case poiUtil.Town:
			matchLen = 6 + 3
		case poiUtil.Village:
			matchLen = 9 + 3
		}

		if len(last.AddressCode) == 12 && len(*paddrCode) == 12 && bytes.Equal(poiUtil.S2B(last.AddressCode)[0:matchLen], poiUtil.S2B(*paddrCode)[0:matchLen]) {
			inRegion = true
		}
	}

	return inRegion
}

// Parse 中文地址要素分割
// 如果当前 level 和 匹配到的 level 相差 2个档次及以上
// 需要根据匹配路径来判别 往哪个方向搜索
// last level = A-B
// cur level = A-B-x-y
// 那么 只能在 A-B 匹配的情况下 搜索，提高准确性
func (as *AddrSegmentation) parseXzqh(bsb *poiUtil.ByteStringBuffer, cam *ChineseAddressModel) bool {
	// 省
	for _, pV := range as.Provinces {
		if as.findXzqh(bsb, &pV) {
			break
		}
	}

	// 市
	if as.hit != -1 && as.pXzqh.Level == poiUtil.Province {
		if citys, ok := as.Citys[as.hit]; ok {
			for _, cV := range citys.Xzqhs {
				if as.findXzqh(bsb, &cV) {
					break
				}
			}
		} else {
			log.Error().Int32("XzqhID", as.hit).Msg("未找到 省 所在城市")
		}
	} else { // 未找到 省, 则扫描全国所有 市
		for _, ciytsV := range as.Citys {
			for _, cV := range ciytsV.Xzqhs {
				if as.findXzqh(bsb, &cV) {
					break
				}
			}
		}
	}

	// 区县
	if as.hit != -1 && as.pXzqh.Level == poiUtil.City {
		if districts, ok := as.Districts[as.hit]; ok {
			for _, dV := range districts.Xzqhs {
				if as.findXzqh(bsb, &dV) {
					break
				}
			}
		} else {
			log.Error().Int32("XzqhID", as.hit).Msg("未找到 市 所在区县")
		}
	} else { // 未找到 区县, 则扫描全国所有 区县
		for _, districtsV := range as.Districts {
			if as.isInRegion(as.pXzqh, &(districtsV.AddressCodeCommon), poiUtil.District) {
				for _, dV := range districtsV.Xzqhs {
					if as.findXzqh(bsb, &dV) {
						break
					}
				}
			}
		}
	}

	// 乡镇街道
	if as.hit != -1 && as.pXzqh.Level == poiUtil.District {
		if towns, ok := as.Towns[as.hit]; ok {
			for _, tV := range towns.Xzqhs {
				if as.findXzqh(bsb, &tV) {
					break
				}
			}
		} else {
			log.Error().Int32("XzqhID", as.hit).Msg("未找到 县 所在区 街道乡镇")
		}
	} else { // 未找到 乡镇 则扫描全国 乡镇
		for _, townsV := range as.Towns {
			if as.isInRegion(as.pXzqh, &(townsV.AddressCodeCommon), poiUtil.Town) {
				for _, tV := range townsV.Xzqhs {
					if as.findXzqh(bsb, &tV) {
						break
					}
				}
			}
		}
	}

	// 乡村，居委会
	if as.hit != -1 && as.pXzqh.Level == poiUtil.Town {
		if villages, ok := as.Villages[as.hit]; ok {
			for _, vV := range villages.Xzqhs {
				if as.findXzqh(bsb, &vV) {
					break
				}
			}
		} else {
			log.Error().Int32("XzqhID", as.hit).Msg("未找到 街道乡镇 所在区 乡村,居委会")
		}
	} else { // 未找到 乡村,居委会 则扫描全国 乡村,居委会
		for _, villagesV := range as.Villages {
			if as.isInRegion(as.pXzqh, &(villagesV.AddressCodeCommon), poiUtil.Village) {
				for _, vV := range villagesV.Xzqhs {
					if as.findXzqh(bsb, &vV) {
						break
					}
				}
			}
		}
	}

	// 根据最后的找到的行政区划条目 填充字段
	as.fillXzqh(cam)

	return true
}

// 解析行政区划之后的地址信息
func (as *AddrSegmentation) parseDetailAddr(bsb *poiUtil.ByteStringBuffer, cam *ChineseAddressModel) bool {
	//fmt.Printf("行政区划分解后的地址=%s", bsb.String())

	// 路
	loc := as.roadRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.Road.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])

		// 路号
		loc = as.roadNumRegex.FindStringIndex(bsb.String())
		if loc != nil {
			cam.RoadNum.WriteString(bsb.String()[loc[0]:loc[1]])
			bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
		}
	}

	// 辅路
	loc = as.roadRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.SubRoad.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])

		// 辅路号
		loc = as.roadNumRegex.FindStringIndex(bsb.String())
		if loc != nil {
			cam.SubRoadNum.WriteString(bsb.String()[loc[0]:loc[1]])
			bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
		}
	}

	// POI
	loc = as.poiRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.Poi.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	// 辅 POI
	loc = as.poiRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.SubPoi.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	// 楼号
	loc = as.buildingNumRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.BuildingNum.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	// 单元号
	loc = as.unitNumRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.UnitNum.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	// 楼层号
	loc = as.floorNumRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.FloorNum.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	// 房号
	loc = as.roomNumRegex.FindStringIndex(bsb.String())
	if loc != nil {
		cam.RoomNum.WriteString(bsb.String()[loc[0]:loc[1]])
		bsb.Replace(loc[0], as.filling, loc[1]-loc[0])
	}

	return true
}

// 处理行政区划和详细地址分割后剩余的地址信息
func (as *AddrSegmentation) clipExtraInfo(bsb *poiUtil.ByteStringBuffer, cam *ChineseAddressModel) {
	continueSep := 0
	needByte := false
	sep := rune(as.filling)

	for _, v := range bsb.String() {
		if v != sep {
			cam.ExtraInfo.WriteRune(v)
			continueSep = 0
			needByte = true
		} else {
			continueSep += 1
			if continueSep == 1 && needByte {
				cam.ExtraInfo.WriteByte(as.extraInfoSep)
			}

			needByte = false
		}
	}

	if cam.ExtraInfo.Len() > 0 && cam.ExtraInfo.String()[cam.ExtraInfo.Len()-1] == as.filling {
		cam.ExtraInfo.Slice(0, cam.ExtraInfo.Len()-1)
	}
}

// Parse 地址分解入口
func (as *AddrSegmentation) Parse(bsb *poiUtil.ByteStringBuffer, cam *ChineseAddressModel) bool {
	as.reset()
	as.parseXzqh(bsb, cam)
	as.parseDetailAddr(bsb, cam)
	as.clipExtraInfo(bsb, cam)
	return true
}

func (as *AddrSegmentation) reset() {
	as.index = -1
	as.pXzqh = nil
	as.hit = -1
}

// Init 分配内存，地址设置填充符
func (as *AddrSegmentation) Init(fillingChar byte, xzqhSep rune, extraInfoSep byte) {
	as.filling = fillingChar
	as.xzqhSep = xzqhSep
	as.extraInfoSep = extraInfoSep
	as.Create()
	as.roadRegex = regexp.MustCompile(`[\p{Han}\dA-Za-z]+(?:路|街[道坊]?|道|组|號|巷|弄|里|亍|胡同|条|坡)+`)
	as.roadNumRegex = regexp.MustCompile(`(?:[东南西北]?[\d一二三四五六七八九十-]+(号院|号)?)`)
	as.buildingNumRegex = regexp.MustCompile(`[\dA-Za-z一二三四五六七八九十-]{1,3}(栋|橦|幢|座|号[楼]?|阁|排)`)
	as.poiRegex = regexp.MustCompile(`[\p{Han}\d]+([小校发新社东南西北园一二三四五六七八九十\dA-Za-z业]区|坞|苑|庭|寓|府邸|[学驾]校|
					[^号]楼|大厦|居|工作室|湾|城|村|乡|園|园|[大中小]学|大?院|(有限)?公司|集团|店|[市广农]场|局|超市|药业|
					中心|企业汇|仓库|鞋业|实业|厂|祠|国际|医院|公馆|府|公社|家纺|服饰|幼儿园|新村|工业园|纺织|宅|宿舍|新村|
					物流|居|服饰|部|山庄|阁|名门|电器|谭|基地|畈|药房|址|企业|镇|港|电商|桥|寺|铝业)`)
	as.unitNumRegex = regexp.MustCompile(`[一二三四五六七八九十甲乙丙0-9]{1,2}(单元|[号]?门|梯|号)`)
	as.floorNumRegex = regexp.MustCompile(`[正负]?[一二三四五六七八九十0-9-]+(层|楼|F)`)
	as.roomNumRegex = regexp.MustCompile(`[\da-zA-Z]{1,4}-?[\da-zA-Z-]{0,2}-?[\da-zA-Z-]{0,4}(室|房|$|号)?`)
}

// findXzqh 查找行政区划条目
// 如果有一条行政区划 "浙江省-金华市-义乌市-江东街道-青口村委会"
// 待分割地址 "浙江-金华-义乌市-青口"
// "青口" 很多个省市有 "青口" 这个名称
// 针对跨级 curLevel - lastLvel > 1 需要 根据省市 判断地址是否一致
// 简单处理，根据zipcode来判定是否属于同一个城市
func (as *AddrSegmentation) findXzqh(bsb *poiUtil.ByteStringBuffer, pXzqhItem *poiUtil.XingZhengQuHuaItem) bool {
	find := false

	if as.index = strings.Index(bsb.String(), pXzqhItem.Name); as.index != -1 {
		bsb.Replace(as.index, as.filling, len(pXzqhItem.Name))
		as.pXzqh = pXzqhItem
		as.hit = pXzqhItem.Id
		find = true
	} else if as.index = strings.Index(bsb.String(), pXzqhItem.ShortName); as.index != -1 {
		bsb.Replace(as.index, as.filling, len(pXzqhItem.ShortName))
		as.pXzqh = pXzqhItem
		as.hit = pXzqhItem.Id
		find = true
	}

	return find
}

// which xzqh level
func (as *AddrSegmentation) whichXzqhLevel(start, end, level int, cam *ChineseAddressModel) {
	if level == poiUtil.Province {
		cam.Province.WriteString(as.pXzqh.MergeName[start:end])
	} else if level == poiUtil.City {
		cam.City.WriteString(as.pXzqh.MergeName[start:end])
	} else if level == poiUtil.District {
		cam.District.WriteString(as.pXzqh.MergeName[start:end])
	} else if level == poiUtil.Town {
		cam.Town.WriteString(as.pXzqh.MergeName[start:end])
	} else if level == poiUtil.Village {
		cam.Village.WriteString(as.pXzqh.MergeName[start:end])
	}
}

// fillXzqh
func (as *AddrSegmentation) fillXzqh(cam *ChineseAddressModel) {
	if as.hit != -1 {
		start := 0
		level := 0
		for i, v := range as.pXzqh.MergeName {
			if v == as.xzqhSep {
				as.whichXzqhLevel(start, i, level, cam)
				start = i + 1 // over '-', one byte
				level += 1
			}
		}

		// for last one
		as.whichXzqhLevel(start, len(as.pXzqh.MergeName), level, cam)
	}

	cam.XzqhCoding.WriteString(as.pXzqh.AddressCode)
	cam.ZipCode.WriteString(as.pXzqh.ZipCode)
	cam.CityCode.WriteString(as.pXzqh.CityCode)
	cam.Lng = as.pXzqh.Lng
	cam.Lat = as.pXzqh.Lat
}
