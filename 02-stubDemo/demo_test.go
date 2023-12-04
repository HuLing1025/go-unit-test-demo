package stubdemo

import (
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	. "github.com/smartystreets/goconvey/convey"
)

// testCaseForInfo test cases for global varible Num.
var testCaseForNum = []struct {
	name string
	to   int
	want int
}{
	{"test-case-01: 测试全局变量 Num 打桩", 99, 99},
}

// testCaseForInfo test cases for global varible Info.
var testCaseForInfo = []struct {
	name string
	to   string
	want string
}{
	{"test-case-01: 测试全局变量 Info 打桩", "This is infomation for you.", "This is infomation for you."},
}

// testCaseForStore test cases for function DoStore.
var testCaseForStore = []struct {
	name string
	to   string
	want string
}{
	{"test-case-01: 测试公有函数 DoStore() 打桩", "do not store", "success"},
	{"test-case-02: 测试公有函数 DoStore() 打桩", "store", "success"},
}

// testCaseForGetInfo test cases for private function getInfo.
var testCaseForGetInfo = []struct {
	name string
	to   string
	want string
}{
	{"测试Mock getInfo 01", "information mock 01", "information mock 01"},
	{"测试Mock getInfo 02", "information mock 02", "information mock 02"},
}

// testCaseForMsg test cases for method GetMsg.
var testCaseForMsg = []struct {
	name string
	to   string
	want string
}{
	{"测试Mock GetMsg 01", "I am init message 01", "This is meaasge from mocking method"},
	{"测试Mock GetMsg 02", "I am init message 02", "This is meaasge from mocking method"},
}

func TestMockVarible(t *testing.T) {
	// test mock Num.
	Convey("测试 Stub global varibles: Num", t, func() {
		for _, v := range testCaseForNum {
			// mock varibles.
			patch := gomonkey.
				ApplyGlobalVar(&Num, v.to)
			defer patch.Reset()

			So(Num, ShouldEqual, v.want)
		}
	})

	// test mock Info.
	Convey("测试 Stub global varibles: Info", t, func() {
		for _, v := range testCaseForInfo {
			// mock varibles.
			patch := gomonkey.
				ApplyGlobalVar(&Info, v.to)
			defer patch.Reset()

			So(Info, ShouldEqual, v.want)
		}
	})
}

func TestMockFunction(t *testing.T) {
	Convey("测试 Stub public fuctions: GetResult()", t, func() {
		for _, v := range testCaseForStore {
			Convey(v.name, func() {
				patch := gomonkey.ApplyFunc(GetResult, func(d Demo) string {
					return "success"
				})
				defer patch.Reset()

				d := Demo{msg: v.to}
				So(GetResult(d), ShouldEqual, v.want)
			})
		}
	})
}

func TestMockPrivateFunction(t *testing.T) {
	Convey("测试 Stub private fuctions: getInfo()", t, func() {
		for _, v := range testCaseForGetInfo {
			Convey(v.name, func() {
				patch := gomonkey.ApplyFunc(getInfo, func() string {
					return v.want
				})
				defer patch.Reset()

				So(getInfo(), ShouldEqual, v.want)
			})
		}
	})
}

func TestMockMethod(t *testing.T) {
	Convey("测试Mock methods", t, func() {
		for _, v := range testCaseForMsg {
			// child convey has no param 't'.
			Convey(v.name, func() {
				d := NewDemo(v.to)
				// mock method. 'd' must be a pointer.
				patch := gomonkey.ApplyMethod(reflect.TypeOf(d), "GetMsg", func(_ *Demo) string {
					return v.want
				})
				defer patch.Reset()

				So(d.GetMsg(), ShouldEqual, v.want)
			})
		}
	})
}
