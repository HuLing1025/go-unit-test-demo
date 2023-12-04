#### What Can Unit Test Provide
- Confidence you can reshape code without worring about changing behavior
- Documentation for humans as to how the system should behave
- Much faster and more reliable feedback than manual testing
  
#### (Well Designed) units
- Easy to write meaningful tests.
- Easy to refactor.
  
#### Small Steps(TDD workflow?)
- Write a small test for a small amount desired behaviour
- Check the test fail with a clear error(red)
- Write the minimal amount of code to make the test pass(green)
- Refactor
- Repeat

#### About Go Test Frameworks
There are too many test frameworks for golang, and we choose `goconvey` and `gomonkey` because of their sample(or neat) and graceful.
- `goconvey`
- `gomonkey`
  
Install `goconvey` web ui.
```sh
go install github.com/smartystreets/goconvey@latest
```

---

#### About Stubbing
Focus on the function under test itself, stubbing all its dependencies(varibles、functions、methods、interfaces). 

Firstly, you need to import the packages:
```go
import (
    "fmt"
    "testing"

    "github.com/prashantv/gostub"
    "bou.ke/monkey"
)
```
- **Stub variables**
  
  1. Define function `GetValue`
  ```go
    func GetInfo() string {
	    return fmt.Sprintf("%d:%s", valueInt, valueStr)
    }
  ```
   2. Use `gostub` to mock the global variables.
  ```go
    // stub variables.
    variableStub := func() *gostub.Stubs {
		return gostub.Stub(&valueInt, 6).Stub(&valueStr, "test_value")
    }

    // Use stubs to mock the variables.
    stub := variableStub()
    // reset stubs after test.
    defer stub.Reset()

    // execute test.
    if got := GetInfo(); got != "6:test_value" {
		t.Errorf("GetInfo() = %v, want %v", got, tt.want)
	}
  ```
  
- **Stub functions**
  1. Define function `DoStore`
  ```go
    func CheckDoStore() string {
        if DoStore() {
            return "Store success."
        }
        return "Store failed."
    }

    func DoStore() bool {
	    return false
    }
  ```
  2. Use `monkey` to mock the functions.
  ```go
    // functions stub.
    monkey.Patch(DoStore, func() bool {
		return true
	})
    // unpatch all.
	defer monkey.UnpatchAll()
    
    want := "Store success."
     // execute test.
    if got := CheckDoStore(); got != want {
		t.Errorf("GetInfo() = %v, want %v", got, want)
	}
  ```

- **Stub methods**
  
  1. Define a struct `MsgModel` which contains a method `GetMsg`.
  ```go
    type MsgModel struct {
        msg   string
    }

    // constructor.
    func NewMsgModel() *MsgModel {
        return &MsgModel{
            msg: "This is init message.",
        }
    }

    // method.
    func (s *MsgModel) GetMsg() string {
        return s.msg
    }

  ```
  2. Use `monkey` to mock the methods.
  ```go
    msgModel := NewMsgModel()

	// 实例方法打桩
	monkey.PatchInstanceMethod(reflect.TypeOf(&msgModel), "GetMsg", func(s *MsgModel) string {
		return "This is mocking message."
	})
	defer monkey.UnpatchAll()

	if got := msgModel.GetMsg(); got != "This is mocking message." {
		t.Errorf("GetMsg() = %v, want %v", got, "This is mocking message.")
	}
  ```

- **Stub interfaces**
  repo/user.go
  ```go
  package repo

  import "gorm.io/gorm"

  type UserModel struct {
  	Id   string `gorm:"column:id" json:"id"`
  	Name string `gorm:"column:name" json:"name"`
  	Age  int    `gorm:"column:age" json:"age"`
  }

  type UserRepo struct {
  	db *gorm.DB
  }

  // UserRepoInterface
  type UserRepoInterface interface {
  	Create(model *UserModel) error
  	SelectOne(id string) (UserModel, error)
  }

  // NewUserRepo
  func NewUserRepo(db *gorm.DB) UserRepoInterface {
  	return &UserRepo{db: db}
  }

  func (u *UserModel) TableName() string {
  	return "user"
  }

  func (r *UserRepo) Create(model *UserModel) error {
  	return r.db.Create(model).Error
  }

  func (r *UserRepo) SelectOne(id string) (user UserModel, err error) {
  	err = r.db.Model(&UserModel{}).
  		Where(`id = ?`, id).
  		First(&user).Error
  	return
  }
  ```
  service/user.go
  ```go
  package service

  import (
    "main/03-mockDemo/consts"
    "main/03-mockDemo/repo"
    "main/03-mockDemo/service/dto"

    "github.com/pkg/errors"
  )

  type UserService struct {
    userRepo repo.UserRepoInterface
    userDto  dto.UserDto
  }

  func NewUserService(userRepo repo.UserRepoInterface, userDto dto.UserDto) *UserService {
    return &UserService{userRepo: userRepo, userDto: userDto}
  }

  func (s *UserService) AddUser() (bool, error) {
    user, err := s.userRepo.SelectOne(s.userDto.Id)
    if err != nil {
      return false, errors.Wrap(err, consts.SystemErrorPrefix)
    }

    // user exist.
    if len(user.Id) != 0 {
      return false, errors.New(consts.UserExist)
    }

    // do create.

    return true, nil
  }
  ```
  service/user_test.go
  ```go
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

  // testCaseForAddUser test cases for interface SelectOne.
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
  			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne",  func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
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
  			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne",  func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
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
  			return gomonkey.ApplyMethod(reflect.TypeOf(&repo.UserRepo{}), "SelectOne",  func(_ *repo.UserRepo, id string) (user repo.UserModel, err error) {
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
  ```

- ~~**Mock SQLs**~~
- ~~**Mock redis**~~
- ~~**Mock http requests**~~
  
---

#### Summary
- If your test case requires stubbing functions, then you need to use the following command to execute the test case, otherwise the function stub will not take effect. Why? Because golang is compiled with inline enabled by default, so short functions are inlined into the function under test when source code is campaigning.
```sh
    go test -gcflags=all=-l .\xxxx_test.go 
```
or
```sh
    go test -gcflags=all=-l
```
- You'll see the web ui of test results after you execute the below command, and browse `http://localhost:8080`.
```sh
    goconvey
```
