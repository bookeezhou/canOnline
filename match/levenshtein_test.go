package match

import (
	"testing"
)

func TestLevenshteinDistance(t *testing.T) {
	var lstn Levenshtein
	lstn.Init()

	str1 := "你我好"
	str2 := "你好123"
	d, r := lstn.Distance(str1, str2)
	t.Logf("%d, %0.2f", d, r)

	str3 := "三个大厦"
	str4 := "三个研发大楼"
	d, r = lstn.Distance(str3, str4)
	t.Logf("%d, %0.2f", d, r)
}
