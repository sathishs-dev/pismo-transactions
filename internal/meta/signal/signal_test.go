package signal

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type signalTestSuite struct {
	suite.Suite
}

func TestSignalSuite(t *testing.T) {
	suite.Run(t, new(signalTestSuite))
}

func (s *signalTestSuite) SetupTest() {
	initGlobalSig()
}

func (s *signalTestSuite) TestDefaultSignalsInterrupt() {
	go func() {
		time.Sleep(time.Millisecond * 20)
		syscall.Kill(os.Getpid(), syscall.SIGINT)
	}()

	var ok bool
	done := New()

	Add(func() {
		ok = true
	})

	<-done
	s.True(ok)
}

func (s *signalTestSuite) TestDefaultSignalsTermination() {
	go func() {
		time.Sleep(time.Millisecond * 20)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	var ok bool
	done := New()

	Add(func() {
		ok = true
	})

	<-done
	s.True(ok)
}

func (s *signalTestSuite) TestCustomSignals() {
	go func() {
		time.Sleep(time.Millisecond * 20)
		syscall.Kill(os.Getpid(), syscall.SIGHUP)
	}()

	var ok bool
	done := New(syscall.SIGHUP)

	Add(func() {
		ok = true
	})

	<-done
	s.True(ok)
}
