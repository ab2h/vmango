// +build unit

package handlers_test

import (
	"fmt"
	"testing"
	"time"
	"vmango/dal"
	"vmango/domain"
	"vmango/testool"

	"github.com/stretchr/testify/suite"
)

const IMAGES_URL = "/images/"
const IMAGES_API_URL = "/api/images/"

type ImageHandlersTestSuite struct {
	suite.Suite
	testool.WebTest
}

func (suite *ImageHandlersTestSuite) TestAuthRequired() {
	rr := suite.DoGet(IMAGES_URL)
	suite.Assert().Equal(302, rr.Code, rr.Body.String())
	suite.Assert().Equal("/login/?next=/images/", rr.Header().Get("Location"))
}

func (suite *ImageHandlersTestSuite) TestAPIAuthRequired() {
	rr := suite.DoGet(IMAGES_API_URL)
	suite.Assert().Equal(401, rr.Code, rr.Body.String())
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.JSONEq(`{"Error": "Authentication failed"}`, rr.Body.String())
}

func (suite *ImageHandlersTestSuite) TestAPIPostNotAllowed() {
	suite.APIAuthenticate("admin", "secret")
	rr := suite.DoPost(IMAGES_API_URL, nil)
	suite.Assert().Equal(501, rr.Code, rr.Body.String())
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
}

func (suite *ImageHandlersTestSuite) TestOk() {
	suite.Authenticate()
	suite.ProviderFactory.Add(&domain.Provider{
		Name: "test1",
		Images: &dal.StubImagerep{Data: []*domain.Image{
			{
				Id:   "test_image.img",
				OS:   "TestOS",
				Arch: domain.ARCH_X86,
				Type: domain.IMAGE_FMT_RAW,
				Date: time.Unix(1484891107, 0),
			},
			{
				Id:   "test_image2.img",
				OS:   "OsTest-4.0",
				Arch: domain.ARCH_X86_64,
				Type: domain.IMAGE_FMT_QCOW2,
				Date: time.Unix(1484831107, 0),
			},
		}},
	})
	suite.ProviderFactory.Add(&domain.Provider{
		Name: "test2",
		Images: &dal.StubImagerep{Data: []*domain.Image{
			{
				Id:   "test_image.img",
				OS:   "TestOS",
				Arch: domain.ARCH_X86,
				Type: domain.IMAGE_FMT_RAW,
				Date: time.Unix(1484891107, 0),
			},
			{
				Id:   "test_image2.img",
				OS:   "OsTest-4.0",
				Arch: domain.ARCH_X86_64,
				Type: domain.IMAGE_FMT_QCOW2,
				Date: time.Unix(1484831107, 0),
			},
		}},
	})

	rr := suite.DoGet(IMAGES_URL)
	suite.Assert().Equal(200, rr.Code, rr.Body.String())
}

func (suite *ImageHandlersTestSuite) TestAPIOk() {
	suite.APIAuthenticate("admin", "secret")
	suite.ProviderFactory.Add(&domain.Provider{
		Name: "test2",
		Images: &dal.StubImagerep{Data: []*domain.Image{
			{
				Id:       "test_image.img",
				OS:       "TestOS",
				Arch:     domain.ARCH_X86,
				Type:     domain.IMAGE_FMT_RAW,
				Date:     time.Unix(1484891107, 0).UTC(),
				PoolName: "hello",
			},
			{
				Id:       "test_image2.img",
				OS:       "OsTest-4.0",
				Arch:     domain.ARCH_X86_64,
				Type:     domain.IMAGE_FMT_QCOW2,
				Date:     time.Unix(1484831107, 0).UTC(),
				PoolName: "hello2",
			},
		}},
	})

	rr := suite.DoGet(IMAGES_API_URL)
	suite.Equal(200, rr.Code, rr.Body.String())
	suite.Equal("application/json; charset=UTF-8", rr.Header().Get("Content-Type"))
	suite.JSONEq(`{
		"Title": "Images",
		"Images": {
			"test2": [{
				"Id": "test_image.img",
				"OS": "TestOS",
				"Arch": "x86",
				"Size": 0,
				"Type": 1,
				"Date": "2017-01-20T05:45:07Z",
				"PoolName": "hello"
			},{
				"Id": "test_image2.img",
				"OS": "OsTest-4.0",
				"Arch": "x86_64",
				"Size": 0,
				"Type": 2,
				"Date": "2017-01-19T13:05:07Z",
				"PoolName": "hello2"
			}]
		}
	}`, rr.Body.String())
}

func (suite *ImageHandlersTestSuite) TestRepFail() {
	suite.Authenticate()
	suite.ProviderFactory.Add(&domain.Provider{
		Name: "test1",
		Images: &dal.StubImagerep{
			ListErr: fmt.Errorf("test repo error"),
		},
	})
	rr := suite.DoGet(IMAGES_URL)
	suite.Assert().Equal(500, rr.Code)
}

func TestImageHandlersTestSuite(t *testing.T) {
	suite.Run(t, new(ImageHandlersTestSuite))
}
