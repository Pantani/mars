package cli_test

import (
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Pantani/mars/testutil/network"
)

type IntegrationTestSuite struct {
	suite.Suite

	locked bool
	net    *network.Network
	cfg    network.Config
}

func (s *IntegrationTestSuite) network() *network.Network {
	s.net = network.New(s.T(), s.cfg)
	s.locked = true
	return s.net
}

func (s *IntegrationTestSuite) waitForNextBlock() {
	s.T().Log("wait for next block")
	s.Require().NoError(s.net.WaitForNextBlock())
}

func (s *IntegrationTestSuite) SetupTest() {
	s.T().Log("setting up test")
	s.cfg = network.DefaultConfig()
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")
	s.cfg = network.DefaultConfig()
}

func (s *IntegrationTestSuite) TearDownTest() {
	s.T().Log("tearing down test")
	if s.net != nil && s.locked {
		s.net.Cleanup()
		s.locked = false
	}
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
