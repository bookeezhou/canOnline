package poiUtil

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strconv"
	"strings"
)

const (
	xzqhElementLen = 12
	Province       = 0 // 省级，直辖市层级
	City           = 1 // 市
	District       = 2 // 区县
	Town           = 3 // 街道乡镇
	Village        = 4 // 村，道路
)

// 行政区划
type XingZhengQuHuaItem struct {
	Id          int32
	Pid         int32   // 父ID
	Level       uint8   // 行政级别 省=0，市=1，区县=2，街道=3，村委会=4
	AddressCode string  // 地址编码
	ZipCode     string  // 邮编
	CityCode    string  // 城市区号
	Name        string  // 当前行政区划名称
	ShortName   string  // 当前行政区划简称
	MergeName   string  // 当前到顶级行政区划合并名称，以’-‘分割
	Pinyin      string  // 当前行政区划名称的拼音
	Lng         float32 // 精度
	Lat         float32 // 维度
}

// xzqhBolocks 同级行政区划集合
type XzqhBlocks struct {
	AddressCodeCommon string // 当前行政区划的上级地址编码
	Xzqhs             []XingZhengQuHuaItem
}

// XingZhengQuHuaDict 行政区划字典
// key=Pid value=Pid所辖区的{市、区县、街道乡镇，村}
// 省级 key=Id, value=省
type XingZhengQuHuaDict struct {
	Provinces map[int32]XingZhengQuHuaItem
	Citys     map[int32]*XzqhBlocks
	Districts map[int32]*XzqhBlocks
	Towns     map[int32]*XzqhBlocks
	Villages  map[int32]*XzqhBlocks
}

// Load 加载行政区划字典
func (xzqhDict *XingZhengQuHuaDict) Load(xzqhFile string) error {
	f, err := os.Open(xzqhFile)
	if err != nil {
		return err
	}
	defer f.Close()

	bufReader := bufio.NewReader(f)
	lineLen := 0
	lineCount := 0
	lineErrorCount := 0
	var xzqhItem XingZhengQuHuaItem
	var bsb ByteStringBuffer

	for {
		line, rError := bufReader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if rError != nil {
			if rError == io.EOF {
				break
			}
			return rError
		}
		lineCount++

		lineSegment := strings.Split(line, "/")
		lineLen = len(lineSegment)
		if lineLen != xzqhElementLen {
			log.Error().Strs("xzqhRecord", lineSegment).Msg("行政区划记录缺少字段")
			lineErrorCount++
			continue
		}

		// copy lines element to struct
		id, _ := strconv.ParseInt(lineSegment[0], 10, 32)
		xzqhItem.Id = int32(id)
		pid, _ := strconv.ParseInt(lineSegment[1], 10, 32)
		xzqhItem.Pid = int32(pid)
		level, _ := strconv.ParseUint(lineSegment[2], 10, 8)
		xzqhItem.Level = uint8(level)
		xzqhItem.AddressCode = lineSegment[3]
		xzqhItem.ZipCode = lineSegment[4]
		xzqhItem.CityCode = lineSegment[5]
		xzqhItem.Name = lineSegment[6]
		xzqhItem.ShortName = lineSegment[7]
		xzqhItem.MergeName = lineSegment[8]
		xzqhItem.Pinyin = lineSegment[9]
		lng, _ := strconv.ParseFloat(lineSegment[10], 32)
		xzqhItem.Lng = float32(lng)
		lat, _ := strconv.ParseFloat(lineSegment[11], 32)
		xzqhItem.Lat = float32(lat)

		if xzqhItem.Level == Province { // 省
			xzqhDict.Provinces[xzqhItem.Id] = xzqhItem
		} else if xzqhItem.Level == City { // 市
			if v, ok := xzqhDict.Citys[xzqhItem.Pid]; ok {
				v.Xzqhs = append(v.Xzqhs, xzqhItem)
			} else {
				// 新的省所辖城市
				xzqhDict.Citys[xzqhItem.Pid] = new(XzqhBlocks)
				xzqhDict.Citys[xzqhItem.Pid].Xzqhs = append(xzqhDict.Citys[xzqhItem.Pid].Xzqhs, xzqhItem)
				bsb.Reset()
				bsb.WriteString(xzqhItem.AddressCode)
				bsb.Replace(2, '0', 2)
				xzqhDict.Citys[xzqhItem.Pid].AddressCodeCommon = bsb.CopyString()
			}
		} else if xzqhItem.Level == District { // 区县
			if v, ok := xzqhDict.Districts[xzqhItem.Pid]; ok {
				v.Xzqhs = append(v.Xzqhs, xzqhItem)
			} else {
				// 新的市所辖区县
				xzqhDict.Districts[xzqhItem.Pid] = new(XzqhBlocks)
				xzqhDict.Districts[xzqhItem.Pid].Xzqhs = append(xzqhDict.Districts[xzqhItem.Pid].Xzqhs, xzqhItem)
				bsb.Reset()
				bsb.WriteString(xzqhItem.AddressCode)
				bsb.Replace(4, '0', 2)
				xzqhDict.Districts[xzqhItem.Pid].AddressCodeCommon = bsb.CopyString()
			}
		} else if xzqhItem.Level == Town { // 街道乡镇
			if v, ok := xzqhDict.Towns[xzqhItem.Pid]; ok {
				v.Xzqhs = append(v.Xzqhs, xzqhItem)
			} else {
				// 新的区县所辖乡镇街道
				xzqhDict.Towns[xzqhItem.Pid] = new(XzqhBlocks)
				xzqhDict.Towns[xzqhItem.Pid].Xzqhs = append(xzqhDict.Towns[xzqhItem.Pid].Xzqhs, xzqhItem)
				bsb.Reset()
				bsb.WriteString(xzqhItem.AddressCode)
				bsb.Replace(6, '0', 3)
				xzqhDict.Towns[xzqhItem.Pid].AddressCodeCommon = bsb.CopyString()
			}
		} else if xzqhItem.Level == Village { // 乡村
			if v, ok := xzqhDict.Villages[xzqhItem.Pid]; ok {
				v.Xzqhs = append(v.Xzqhs, xzqhItem)
			} else {
				// 新的乡镇所辖乡村
				xzqhDict.Villages[xzqhItem.Pid] = new(XzqhBlocks)
				xzqhDict.Villages[xzqhItem.Pid].Xzqhs = append(xzqhDict.Villages[xzqhItem.Pid].Xzqhs, xzqhItem)
				bsb.Reset()
				bsb.WriteString(xzqhItem.AddressCode)
				bsb.Replace(9, '0', 3)
				xzqhDict.Villages[xzqhItem.Pid].AddressCodeCommon = bsb.CopyString()
			}
		}
	}

	log.Info().Int("xzqhs", lineCount).Int("xzqhsError", lineErrorCount).Msg("行政区划加载情况")
	return nil
}

// Create 初始化
func (xzqhDict *XingZhengQuHuaDict) Create() {
	xzqhDict.Provinces = make(map[int32]XingZhengQuHuaItem)
	xzqhDict.Citys = make(map[int32]*XzqhBlocks)
	xzqhDict.Districts = make(map[int32]*XzqhBlocks)
	xzqhDict.Towns = make(map[int32]*XzqhBlocks)
	xzqhDict.Villages = make(map[int32]*XzqhBlocks)
}
