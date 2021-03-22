package match

import "unicode/utf8"

type Levenshtein struct {
	f       []int
	maxSize int
}

// 设置最大内存，防止频繁申请分配内存，最大100个字符
func (lstn *Levenshtein) Init() {
	lstn.maxSize = 100
	lstn.f = make([]int, lstn.maxSize+1)
}

func (lstn *Levenshtein) reset() {
	for j := range lstn.f {
		lstn.f[j] = j
	}
}

// -1: 表示待分析字符串 b 字符个数长于 maxSize
// 返回: (距离，距离相似度)
func (lstn *Levenshtein) Distance(a, b string) (int, float64) {
	bCount := utf8.RuneCountInString(b)
	if bCount+1 > lstn.maxSize {
		return -1, 0.0
	}

	lstn.reset()

	for _, ca := range a {
		j := 1
		fj1 := lstn.f[0] // fj1 is the value of f[j - 1] in last iteration
		lstn.f[0]++
		for _, cb := range b {
			mn := min(lstn.f[j]+1, lstn.f[j-1]+1) // delete & insert
			if cb != ca {
				mn = min(mn, fj1+1) // change
			} else {
				mn = min(mn, fj1) // matched
			}

			fj1, lstn.f[j] = lstn.f[j], mn // save f[j] to fj1(j is about to increase), update f[j] to mn
			j++
		}
	}

	aCount := utf8.RuneCountInString(a)
	ratio := 1 - float64(lstn.f[bCount])/float64(max(aCount, bCount))
	return lstn.f[bCount], ratio
}

func min(a, b int) int {
	if a <= b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}

	return b
}
