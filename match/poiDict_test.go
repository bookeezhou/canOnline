package match

import (
	"path"
	"path/filepath"
	"testing"
)

func TestPoiLoad(t *testing.T) {

	cw, _ := filepath.Abs(filepath.Dir("D:\\WorkFree\\goPrjs\\"))

	err := GetPOIS().Load(filepath.FromSlash(path.Join(cw, "data/zj_poi.tsv")))
	if err != nil {
		t.Errorf("打开文件错误：%s", err)
	}
}
