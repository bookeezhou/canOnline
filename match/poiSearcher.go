package match

import (
	"canOnline/segment"
)

type PoiSearcher struct {
	xzqhCodeCnt       int          // poi字典表 key 长度,默认设定 6
	lvst              *Levenshtein // 编辑距离算法
	distancePercision float64      // 编辑距离精度
}

func (ps *PoiSearcher) Init() {
	ps.xzqhCodeCnt = 6
	ps.lvst = new(Levenshtein)
	ps.lvst.Init()
	ps.distancePercision = 0.85
}

// Math
func (ps *PoiSearcher) Match(cam *segment.ChineseAddressModel) {
	if cam.XzqhCoding.Len() < ps.xzqhCodeCnt {
		return
	}

	if pois, ok := GetPOIS().Pois[cam.XzqhCoding.String()[0:ps.xzqhCodeCnt]]; ok {
		for _, poi := range pois {
			if _, r := ps.lvst.Distance(cam.Poi.String(), poi.Name); r < ps.distancePercision {
				if _, r = ps.lvst.Distance(cam.SubPoi.String(), poi.Name); r < ps.distancePercision {
					// TODO: 路 + 路号 没做匹配
					continue
				}
			}

			ps.writePoiInfo(cam, &poi)
			break
		}
	}
}

func (ps *PoiSearcher) writePoiInfo(cam *segment.ChineseAddressModel, p *Poi) {
	cam.Telephone.WriteString(p.TelePhone)
	cam.PoiType.WriteString(p.Type)
	cam.Lng = p.WgsLng
	cam.Lat = p.WgsLat
}
