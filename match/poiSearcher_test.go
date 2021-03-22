package match

import (
	"canOnline/poiUtil"
	"canOnline/segment"
	"path"
	"path/filepath"
	"testing"
)

func TestPoiMatch(t *testing.T) {
	cw, _ := filepath.Abs(filepath.Dir("D:\\WorkFree\\goPrjs\\"))

	err := GetPOIS().Load(filepath.FromSlash(path.Join(cw, "data/zj_poi.tsv")))
	if err != nil {
		t.Errorf("打开文件错误：%s", err)
	}

	var addrSeg segment.AddrSegmentation
	var bsb poiUtil.ByteStringBuffer
	bsb.Grow(300)
	var chAM segment.ChineseAddressModel
	chAM.Init('-')
	var ps PoiSearcher
	ps.Init()

	addrSeg.Init('&', '-', ',')
	err = addrSeg.Load("D:\\WorkFree\\goPrjs\\data\\nac.csv")
	if err != nil {
		t.Errorf("打开文件错误：%s", err)
	}

	bsb.WriteString("浙江台州黄岩区西城街道西街小区南区7栋二单元4楼403室")
	t.Logf("\n%s\n", bsb.String())
	addrSeg.Parse(&bsb, &chAM)
	ps.Match(&chAM)
	chAM.Format()
	t.Logf("%s", chAM.AddrFormatCache.String())

	bsb.Reset()
	chAM.Clean()
	bsb.WriteString("浙江金华义乌市青口")
	t.Logf("\n%s\n", bsb.String())
	addrSeg.Parse(&bsb, &chAM)
	ps.Match(&chAM)
	chAM.Format()
	t.Logf("%s", chAM.AddrFormatCache.String())

	bsb.Reset()
	chAM.Clean()
	bsb.WriteString("浙江宁波宁海县桃源街道 金山路35号（金山国际西侧）6楼")
	t.Logf("\n%s\n", bsb.String())
	addrSeg.Parse(&bsb, &chAM)
	ps.Match(&chAM)
	chAM.Format()
	t.Logf("%s", chAM.AddrFormatCache.String())

	bsb.Reset()
	chAM.Clean()
	bsb.WriteString("浙江省金华市兰溪市泰苑宾馆对面7-8号")
	t.Logf("\n%s\n", bsb.String())
	addrSeg.Parse(&bsb, &chAM)
	ps.Match(&chAM)
	chAM.Format()
	t.Logf("%s", chAM.AddrFormatCache.String())
}
