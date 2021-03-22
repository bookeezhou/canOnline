package poiUtil

import "testing"

func TestXzqhLoad(t *testing.T) {
	var xzqhDict XingZhengQuHuaDict
	xzqhDict.Create()

	err := xzqhDict.Load("D:\\WorkFree\\goPrjs\\data\\nac.csv")
	if err != nil {
		t.Errorf("打开文件错误：%s", err)
	}

	cnt := 0
	cnt = len(xzqhDict.Provinces)
	if cnt != 34 {
		t.Errorf("正确省个数=%d, 当前省个数=%d", 34, cnt)
	}

	cnt = len(xzqhDict.Citys)
	if cnt != 34 {
		t.Errorf("正确市所属的省个数=%d, 当前市所属的省个数=%d", 34, cnt)
	}
	for _, v := range xzqhDict.Citys {
		t.Logf("%s\n", v.AddressCodeCommon)
	}

	cnt = len(xzqhDict.Districts)
	if cnt != 373 {
		t.Errorf("正确区县所属的市个数=%d, 当前正确区县所属的市个数=%d", 373, cnt)
	}
	for _, v := range xzqhDict.Districts {
		t.Logf("区县%s\n", v.AddressCodeCommon)
	}

	cnt = len(xzqhDict.Towns)
	if cnt != 3232 {
		t.Errorf("正确乡镇街道所属的区县个数=%d, 当前乡镇街道所属的区县个数=%d", 3232, cnt)
	}

	cnt = len(xzqhDict.Villages)
	if cnt != 42868 {
		t.Errorf("正确村所属的乡镇街道个数=%d, 当前村所属的乡镇街道个数=%d", 42868, cnt)
	}

	cnt = 0
	for _, v := range xzqhDict.Villages {
		cnt += len(v.Xzqhs)
	}
	if cnt != 666655 {
		t.Errorf("正确的乡村个数=%d, 当前乡村个数=%d", 666655, cnt)
	}
}
