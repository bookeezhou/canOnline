package poiUtil

import (
	//"bytes"
	"fmt"
	"strings"
	"testing"
)

func BenchmarkBSBffer(b *testing.B) {
	var bsb ByteStringBuffer
	for i := 0; i < b.N; i++ {
		//fmt.Fprint(&bsb, "aaa")
		bsb.WriteString("aaa")
		_ = bsb.String()
		bsb.Reset()
	}

	//fmt.Printf("%s=%d\n", bsb.String(), bsb.Len())
}

func BenchmarkBuilder(b *testing.B) {
	var bsb strings.Builder
	for i := 0; i < b.N; i++ {
		//fmt.Fprint(&bsb, "aaa")
		bsb.WriteString("aaa")
		_ = bsb.String()
		bsb.Reset()
	}

	//fmt.Printf("%s=%d\n", bsb.String(), bsb.Len())
}

func TestReplace(t *testing.T) {
	var bsb ByteStringBuffer
	bsb.WriteString("浙江金华义乌市稠江街道龙回三区1栋1单元303")
	city := "金华"
	index := strings.Index(bsb.String(), city)
	_, err := bsb.Replace(index-10, '&', len(city))
	fmt.Println(err)
	fmt.Println(bsb.String())
}

func TestString(t *testing.T) {
	var bsb ByteStringBuffer
	var a, b string

	bsb.WriteString("浙江省宁波市想上线")
	t.Logf("%s\n", bsb.String())

	t.Logf("%d=", bsb.Len())
	//bsb.Slice(0, bsb.Len())
	a = bsb.CopyString()

	t.Logf("a=%s\n", a)

	bsb.Reset()
	bsb.WriteString("浙xxxx江")
	t.Logf("%d=", bsb.Len())

	b = bsb.CopyString()
	t.Logf("b=%s\n", b)
	t.Logf("a=%s\n", a)
	t.Logf("==>%s\n", bsb.String())
}
