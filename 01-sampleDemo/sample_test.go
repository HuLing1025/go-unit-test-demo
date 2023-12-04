package sampledemo

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testCases = []struct {
	name   string
	a      string
	b      string
	want   string
	assert Assertion
}{
	{
		name:   "测试用例-01",
		a:      "hello",
		b:      "goconvey",
		want:   "hello goconvey",
		assert: ShouldEqual,
	},
	{
		name:   "测试用例-02",
		a:      "welcome",
		b:      "to China",
		want:   "welcome to China",
		assert: ShouldEqual,
	},
	{
		name:   "测试用例-03",
		a:      "test",
		b:      "cases",
		want:   "test cases",
		assert: ShouldEqual,
	},
}

func TestStrMerage(t *testing.T) {
	Convey("测试字符串合并", t, func() {
		for _, tc := range testCases {
			Convey(tc.name, func() {
				So(StrMerage(tc.a, tc.b), tc.assert, tc.want)
			})
		}
	})
}
