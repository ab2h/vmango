// +build unit

package handlers_test

import (
	"bytes"
	"net/url"
	"testing"
	"vmango/cfg"
	"vmango/dal"
	"vmango/domain"
	"vmango/testool"

	"github.com/stretchr/testify/suite"
)

const CREATE_URL = "/machines/add/"
const CREATE_API_URL = "/api/machines/"

type MachineCreateHandlerTestSuite struct {
	suite.Suite
	testool.WebTest
	Machines *dal.StubMachinerep
}

func (suite *MachineCreateHandlerTestSuite) SetupTest() {
	suite.WebTest.SetupTest()
	sshkeys := dal.NewConfigSSHKeyrep([]cfg.SSHKeyConfig{
		{Name: "first", Public: "hello"},
	})
	*suite.SSHKeys = *sshkeys
	plans := dal.NewConfigPlanrep([]cfg.PlanConfig{
		{Name: "test-1", Memory: 512 * 1024 * 1024, Cpus: 1, DiskSize: 5},
		{Name: "test-2", Memory: 1024 * 1024 * 1024, Cpus: 2, DiskSize: 10},
	})
	*suite.Plans = *plans
	suite.Machines = &dal.StubMachinerep{}
	suite.ProviderFactory.Add(&domain.Provider{
		Name:     "test1",
		Machines: suite.Machines,
		Images: &dal.StubImagerep{
			Data: []*domain.Image{
				{OS: "TestOS-1.0", Arch: domain.ARCH_X86_64, Size: 10 * 1024 * 1024, Type: domain.IMAGE_FMT_QCOW2, Id: "TestOS-1.0_amd64.img", PoolName: "test"},
			},
		},
	})
	suite.ProviderFactory.Add(&domain.Provider{
		Name:     "test2",
		Machines: suite.Machines,
		Images: &dal.StubImagerep{
			Data: []*domain.Image{
				{OS: "TestOS-1.0", Arch: domain.ARCH_X86_64, Size: 10 * 1024 * 1024, Type: domain.IMAGE_FMT_QCOW2, Id: "TestOS-1.0_amd64.img", PoolName: "test"},
			},
		},
	})
}

func (suite *MachineCreateHandlerTestSuite) TestGetAuthRequired() {
	rr := suite.DoGet(CREATE_URL)
	suite.Equal(302, rr.Code, rr.Body.String())
	suite.Equal(rr.Header().Get("Location"), "/login/?next="+CREATE_URL)
}

func (suite *MachineCreateHandlerTestSuite) TestPostAuthRequired() {
	rr := suite.DoPost(CREATE_URL, nil)
	suite.Equal(302, rr.Code, rr.Body.String())
	suite.Equal(rr.Header().Get("Location"), "/login/?next="+CREATE_URL)
}

func (suite *MachineCreateHandlerTestSuite) TestPostAPIAuthRequired() {
	rr := suite.DoPost(CREATE_API_URL, nil)
	suite.Equal(401, rr.Code, rr.Body.String())
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.JSONEq(`{"Error": "Authentication failed"}`, rr.Body.String())
}

func (suite *MachineCreateHandlerTestSuite) TestBadHTTPMethodNotAllowed() {
	suite.Authenticate()
	rr := suite.DoBad(CREATE_URL)
	suite.Equal(501, rr.Code)
}

func (suite *MachineCreateHandlerTestSuite) TestBadHTTPMethodAPINotAllowed() {
	suite.APIAuthenticate("admin", "secret")
	rr := suite.DoBad(CREATE_API_URL)
	suite.Equal(501, rr.Code, rr.Body.String())
}

func (suite *MachineCreateHandlerTestSuite) TestGetOk() {
	suite.Authenticate()
	rr := suite.DoGet(CREATE_URL)
	suite.Equal("text/html; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.Equal(200, rr.Code, rr.Body.String())
}

func (suite *MachineCreateHandlerTestSuite) TestCreateOk() {
	suite.Authenticate()
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"test-1"},
		"Image":    []string{"TestOS-1.0_amd64.img"},
		"SSHKey":   []string{"first"},
		"Provider": []string{"test1"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_URL, data)
	suite.Equal(302, rr.Code, rr.Body.String())
	suite.Equal(DETAIL_URL("test1", "stub-machine-id"), rr.Header().Get("Location"))
}

func (suite *MachineCreateHandlerTestSuite) TestCreateAPIOk() {
	suite.APIAuthenticate("admin", "secret")
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"test-1"},
		"Image":    []string{"TestOS-1.0_amd64.img"},
		"SSHKey":   []string{"first"},
		"Provider": []string{"test1"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_API_URL, data)
	suite.Equal(201, rr.Code, rr.Body.String())
	suite.Equal(DETAIL_API_URL("test1", "stub-machine-id"), rr.Header().Get("Location"))
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.JSONEq(`{"Message": "Machine testvm (stub-machine-id) created"}`, rr.Body.String())
}

func (suite *MachineCreateHandlerTestSuite) TestCreateNoPlanFail() {
	suite.Authenticate()
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"doesntexist"},
		"Image":    []string{"TestOS-1.0_amd64.img"},
		"SSHKey":   []string{"first"},
		"Provider": []string{"test2"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_URL, data)
	suite.Equal(500, rr.Code)
	suite.Contains(rr.Body.String(), "plan &#34;doesntexist&#34; not found")
	suite.Equal(rr.Header().Get("Location"), "")
}

func (suite *MachineCreateHandlerTestSuite) TestCreateNoImageFail() {
	suite.Authenticate()
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"test-1"},
		"Image":    []string{"doesntexist"},
		"SSHKey":   []string{"first"},
		"Provider": []string{"test1"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_URL, data)
	suite.Equal(500, rr.Code)
	suite.Contains(rr.Body.String(), "image &#34;doesntexist&#34; not found")
	suite.Equal(rr.Header().Get("Location"), "")
}

func (suite *MachineCreateHandlerTestSuite) TestCreateNoHypervisorFail() {
	suite.Authenticate()
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"test-1"},
		"Image":    []string{"TestOS-1.0_amd64.img"},
		"SSHKey":   []string{"first"},
		"Provider": []string{"doesntexist"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_URL, data)
	suite.Equal(500, rr.Code)
	suite.Contains(rr.Body.String(), "provider &#34;doesntexist&#34;: not found")
	suite.Equal(rr.Header().Get("Location"), "")
}

func (suite *MachineCreateHandlerTestSuite) TestCreateNoSSHKeyFail() {
	suite.Authenticate()
	data := bytes.NewBufferString((url.Values{
		"Name":     []string{"testvm"},
		"Plan":     []string{"test-1"},
		"Image":    []string{"TestOS-1.0_amd64.img"},
		"SSHKey":   []string{"doesntexist"},
		"Provider": []string{"test1"},
	}).Encode())
	suite.T().Log(data)
	rr := suite.DoPost(CREATE_URL, data)
	suite.Equal(500, rr.Code)
	suite.Contains(rr.Body.String(), "ssh key &#39;doesntexist&#39; doesn&#39;t exist")
	suite.Equal(rr.Header().Get("Location"), "")
}

func TestMachineCreateHandlerTestSuite(t *testing.T) {
	suite.Run(t, new(MachineCreateHandlerTestSuite))
}
