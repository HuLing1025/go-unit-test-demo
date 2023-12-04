package service

import (
	"errors"
	"main/03-mockDemo/consts"
	"main/03-mockDemo/repo"
	"main/03-mockDemo/service/dto"
	"reflect"
	"testing"

	"github.com/agiledragon/gomonkey"
	. "github.com/smartystreets/goconvey/convey"
)

// testCaseForAddUser test cases for private function getInfo.
var testCaseForAddUser = []struct {
	name    string
	userDto dto.UserDto
	patch   func() *gomonkey.Patches
	result  bool
	err     error
}{
	{
		name: "测试-01 测试系统错误",
		userDto: dto.UserDto{
			Id:   "123",
			Name: "test",
			Age:  18,
		},
		patch: func() *gomonkey.Patches {
			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne", func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
				err = errors.New("数据库内部错误")
				return
			})
		},
		result: false,
		err:    errors.New(consts.SystemErrorPrefix),
	},
	{
		name: "测试-02 用户不存在, 新增成功",
		userDto: dto.UserDto{
			Id:   "123",
			Name: "test",
			Age:  18,
		},
		patch: func() *gomonkey.Patches {
			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne", func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
				return
			})
		},
		result: true,
		err:    nil,
	},
	{
		name: "测试-03 用户存在, 新增失败",
		userDto: dto.UserDto{
			Id:   "123",
			Name: "test",
			Age:  18,
		},
		patch: func() *gomonkey.Patches {
			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne", func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
				user.Id = id
				return
			})
		},
		result: false,
		err:    errors.New(consts.UserExist),
	},
}

func TestAddUser(t *testing.T) {
	Convey("测试 AddUser", t, func() {
		for _, v := range testCaseForAddUser {
			// child convey has no param 't'.
			Convey(v.name, func() {
				s := NewUserService(repo.NewUserRepo(nil), v.userDto)
				patch := v.patch()
				defer patch.Reset()

				// begin test.
				result, err := s.AddUser()

				So(result, ShouldEqual, v.result)
				if err == nil || v.err == nil {
					So(err, ShouldEqual, v.err)
				} else {
					So(err.Error(), ShouldContainSubstring, v.err.Error())
				}
			})
		}
	})
}
