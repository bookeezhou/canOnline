package match

import (
	"bufio"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Poi struct {
	Name      string // poi name
	Address   string // poi 详细地址
	TelePhone string
	Type      string // 地址类型
	AreaCode  string // 全国行政区划编码，省市区三级
	WgsLng    float32
	WgsLat    float32
}

type poiDict struct {
	Pois map[string][]Poi
}

func (pd *poiDict) Init() {
	pd.Pois = make(map[string][]Poi)
}

// 加载POI数据
func (pd *poiDict) Load(poiPath string) error {
	f, err := os.Open(poiPath)
	if err != nil {
		return err
	}
	defer f.Close()

	fixColumns := 13
	bufReader := bufio.NewReader(f)
	columns := 0
	lines := 0
	errorLines := 0

	for {
		line, rError := bufReader.ReadString('\n')
		line = strings.TrimRight(line, "\n")
		if rError != nil {
			if rError == io.EOF {
				break
			}
			return rError
		}
		lines++

		// 跳过首行，行首是列名称
		if lines == 1 {
			continue
		}

		lineSeg := strings.Split(line, "\t")
		columns = len(lineSeg)
		if columns != fixColumns {
			log.Error().Strs("PoiRecord", lineSeg).Msg("POI记录缺少字段")
			errorLines++
			continue
		}

		var poi Poi
		poi.Name = lineSeg[3]
		poi.Address = lineSeg[4]
		poi.TelePhone = lineSeg[5]
		poi.Type = lineSeg[6]
		poi.AreaCode = lineSeg[7]
		lng, _ := strconv.ParseFloat(lineSeg[8], 32)
		poi.WgsLng = float32(lng)
		lat, _ := strconv.ParseFloat(lineSeg[9], 32)
		poi.WgsLat = float32(lat)

		if pois, ok := pd.Pois[poi.AreaCode]; ok {
			pois = append(pois, poi)
		} else {
			pd.Pois[poi.AreaCode] = append(pd.Pois[poi.AreaCode], poi)
		}
	}

	log.Info().Int("PoiRecord", lines).Int("PoiRecord", errorLines).Msg("poi字典加载完毕")
	return nil
}

var poiSingleton *poiDict
var once sync.Once

func GetPOIS() *poiDict {
	once.Do(func() {
		poiSingleton = &poiDict{}
		poiSingleton.Init()
	})

	return poiSingleton
}
