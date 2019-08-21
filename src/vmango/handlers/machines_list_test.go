// +build unit

package handlers_test

import (
	"fmt"
	"testing"
	"vmango/dal"
	"vmango/domain"
	"vmango/testool"

	"github.com/stretchr/testify/suite"
)

const LIST_URL = "/machines/"
const LIST_API_URL = "/api/machines/"

type MachineListHandlerTestSuite struct {
	suite.Suite
	testool.WebTest
	Repo *dal.StubMachinerep
}

func (suite *MachineListHandlerTestSuite) SetupTest() {
	suite.WebTest.SetupTest()
	suite.Repo = &dal.StubMachinerep{}
	suite.ProviderFactory.Add(&domain.Provider{
		Name:     "test1",
		Machines: suite.Repo,
	})
}

func (suite *MachineListHandlerTestSuite) TestAuthRequired() {
	rr := suite.DoGet(LIST_URL)
	suite.Equal(302, rr.Code, rr.Body.String())
	suite.Equal(rr.Header().Get("Location"), "/login/?next="+LIST_URL)
}

func (suite *MachineListHandlerTestSuite) TestAPIAuthRequired() {
	rr := suite.DoGet(LIST_API_URL)
	suite.Equal(401, rr.Code, rr.Body.String())
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.JSONEq(`{"Error": "Authentication failed"}`, rr.Body.String())
}

func (suite *MachineListHandlerTestSuite) TestHTMLOk() {
	suite.Authenticate()
	suite.Repo.ListResponse.Machines = &domain.VirtualMachineList{}
	suite.Repo.ListResponse.Machines.Add(&domain.VirtualMachine{
		Id:   "deadbeefdeadbeefdeadbeefdeadbeef",
		Name: "test",
	})
	rr := suite.DoGet(LIST_URL)
	suite.Equal(200, rr.Code, rr.Body.String())
	suite.Equal("text/html; charset=UTF-8", rr.Header().Get("Content-Type"))
}

func (suite *MachineListHandlerTestSuite) TestJSONOk() {
	suite.APIAuthenticate("admin", "secret")
	suite.Repo.ListResponse.Machines = &domain.VirtualMachineList{}
	suite.Repo.ListResponse.Machines.Add(&domain.VirtualMachine{
		Name:    "test",
		Id:      "123uuid",
		Memory:  456,
		Cpus:    1,
		HWAddr:  "hw:hw:hw",
		VNCAddr: "vnc",
		Creator: "hello",
		ImageId: "stub-image",
		Plan:    "xxxzzz",
		OS:      "WoW",
		Arch:    domain.ARCH_UNKNOWN,
		Ip: &domain.IP{
			Address: "1.1.1.1",
		},
		RootDisk: &domain.VirtualMachineDisk{
			Size:   123,
			Driver: "hello",
			Type:   "wow",
		},
		SSHKeys: []*domain.SSHKey{
			{Name: "test", Public: "keykeykey"},
		},
	})
	suite.Repo.ListResponse.Machines.Add(&domain.VirtualMachine{
		Name:    "hello",
		Id:      "321uuid",
		Memory:  67897,
		Cpus:    4,
		HWAddr:  "xx:xx:xx",
		VNCAddr: "VVV",
		Creator: "wow",
		Plan:    "testx",
		ImageId: "stub-image",
		Ip: &domain.IP{
			Address: "2.2.2.2",
		},
		RootDisk: &domain.VirtualMachineDisk{
			Size:   321,
			Driver: "ehlo",
			Type:   "www",
		},
		SSHKeys: []*domain.SSHKey{
			{Name: "test2", Public: "kekkekkek"},
		},
	})

	rr := suite.DoGet(LIST_API_URL)
	suite.Require().Equal(200, rr.Code, rr.Body.String())
	suite.Require().Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	expected := `{
      "Title": "Machines",
      "Machines": {
        "test1": [{
          "Id": "123uuid",
          "Name": "test",
          "Memory": 456,
          "Cpus": 1,
          "Creator": "hello",
          "Plan": "xxxzzz",
          "Ip": {"Address": "1.1.1.1", "Gateway": "", "Netmask": 0, "UsedBy": ""},
          "HWAddr": "hw:hw:hw",
          "VNCAddr": "vnc",
          "ImageId": "stub-image",
          "OS": "WoW",
          "Arch": "unknown",
          "RootDisk": {
            "Size": 123,
            "Driver": "hello",
            "Type": "wow"
           },
          "SSHKeys": [
            {"Name": "test", "Public": "keykeykey"}
          ]
        }, {
          "Id": "321uuid",
          "Name": "hello",
          "Memory": 67897,
          "Cpus": 4,
          "Creator": "wow",
          "Plan": "testx",
          "OS": "",
          "Arch": "unknown",
          "HWAddr": "xx:xx:xx",
          "VNCAddr": "VVV",
          "ImageId": "stub-image",
          "Ip": {"Address": "2.2.2.2", "Gateway": "", "Netmask": 0, "UsedBy": ""},
          "RootDisk": {
            "Size": 321,
            "Driver": "ehlo",
            "Type": "www"
           },
          "SSHKeys": [
            {"Name": "test2", "Public": "kekkekkek"}
          ]
        }]
      }

    }`
	suite.JSONEq(expected, rr.Body.String())
}

func (suite *MachineListHandlerTestSuite) TestRepFail() {
	suite.Authenticate()
	suite.Repo.ListResponse.Error = fmt.Errorf("test error")
	rr := suite.DoGet(LIST_URL)
	suite.Equal(500, rr.Code, rr.Body.String())
}

func TestMachineListHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MachineListHandlerTestSuite))
}
