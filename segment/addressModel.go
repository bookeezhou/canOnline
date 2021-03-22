package segment

import (
	"canOnline/poiUtil"
	"strconv"
)

// ChineseAddressModel 存放分割后的各个地址要素
type ChineseAddressModel struct {
	Province        poiUtil.ByteStringBuffer // 省(直辖市)
	City            poiUtil.ByteStringBuffer // 市
	District        poiUtil.ByteStringBuffer // 区县
	Town            poiUtil.ByteStringBuffer // 街道(乡镇)
	Village         poiUtil.ByteStringBuffer // 社区(乡村)
	XzqhCoding      poiUtil.ByteStringBuffer // 国家行政区划编码，12位字符
	Road            poiUtil.ByteStringBuffer // 路
	RoadNum         poiUtil.ByteStringBuffer // 路牌号
	SubRoad         poiUtil.ByteStringBuffer // 辅路
	SubRoadNum      poiUtil.ByteStringBuffer // 辅路牌号
	Poi             poiUtil.ByteStringBuffer // 第一兴趣点，小区也是兴趣点
	SubPoi          poiUtil.ByteStringBuffer // 第二兴趣点
	BuildingNum     poiUtil.ByteStringBuffer // 楼号
	UnitNum         poiUtil.ByteStringBuffer // 单元号
	FloorNum        poiUtil.ByteStringBuffer // 楼层号
	RoomNum         poiUtil.ByteStringBuffer // 房号
	Position        poiUtil.ByteStringBuffer // 方位
	ZipCode         poiUtil.ByteStringBuffer // 邮编
	CityCode        poiUtil.ByteStringBuffer // 区号
	Lng             float32                  // 精度
	Lat             float32                  // 维度
	ExtraInfo       poiUtil.ByteStringBuffer // 地址无法分割后的剩余地址信息
	AddrFormatCache poiUtil.ByteStringBuffer // 地址格式化缓存
	Separator       byte                     // 地址要素分割符
	Telephone       poiUtil.ByteStringBuffer // 联系电话
	PoiType         poiUtil.ByteStringBuffer // 地址类型
}

// Format 格式化地址要素,默认'$'分割字段
func (cam *ChineseAddressModel) Format() {
	cam.AddrFormatCache.WriteString(cam.Province.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.City.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.District.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Town.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Village.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.XzqhCoding.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Road.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.RoadNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.SubRoad.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.SubRoadNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Poi.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.SubPoi.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.BuildingNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.UnitNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.FloorNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.RoomNum.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Position.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.ZipCode.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.CityCode.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.ExtraInfo.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(strconv.FormatFloat(float64(cam.Lng), 'f', 6, 32))
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(strconv.FormatFloat(float64(cam.Lat), 'f', 6, 32))
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.Telephone.String())
	cam.AddrFormatCache.WriteByte(cam.Separator)

	cam.AddrFormatCache.WriteString(cam.PoiType.String())
}

// Clean 清除各字段值
func (cam *ChineseAddressModel) Clean() {
	cam.Province.Reset()
	cam.City.Reset()
	cam.District.Reset()
	cam.Town.Reset()
	cam.Village.Reset()
	cam.XzqhCoding.Reset()
	cam.Poi.Reset()
	cam.SubPoi.Reset()
	cam.Road.Reset()
	cam.RoadNum.Reset()
	cam.SubRoad.Reset()
	cam.SubRoadNum.Reset()
	cam.BuildingNum.Reset()
	cam.UnitNum.Reset()
	cam.FloorNum.Reset()
	cam.RoomNum.Reset()
	cam.Position.Reset()
	cam.ZipCode.Reset()
	cam.CityCode.Reset()
	cam.Lng = 0
	cam.Lat = 0
	cam.ExtraInfo.Reset()
	cam.AddrFormatCache.Reset()
	cam.Telephone.Reset()
	cam.PoiType.Reset()
}

// Init 初始化,设置分隔符
func (cam *ChineseAddressModel) Init(separator byte) {
	cam.Separator = separator
	cam.Province.Grow(45) // 15中文字符
	cam.City.Grow(45)
	cam.District.Grow(45)
	cam.Town.Grow(45)
	cam.Village.Grow(45)
	cam.XzqhCoding.Grow(15) // 12个数字
	cam.Poi.Grow(150)
	cam.SubPoi.Grow(150)
	cam.Road.Grow(150)
	cam.RoadNum.Grow(30)
	cam.SubRoad.Grow(150)
	cam.SubRoadNum.Grow(30)
	cam.BuildingNum.Grow(30)
	cam.UnitNum.Grow(30)
	cam.FloorNum.Grow(30)
	cam.RoomNum.Grow(30)
	cam.Position.Grow(30)
	cam.ZipCode.Grow(10)
	cam.City.Grow(10)
	cam.ExtraInfo.Grow(150) // 150字节,50 个中文字符(UTF8编码)
	cam.AddrFormatCache.Grow(600)
	cam.Telephone.Grow(50)
	cam.PoiType.Grow(150)
}
