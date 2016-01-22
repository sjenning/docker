package main

import (
	"bytes"
	"os/exec"
	"strings"

	"github.com/docker/docker/pkg/integration/checker"
	"github.com/docker/docker/registry"
	"github.com/go-check/check"
)

func (s *DockerSuite) TestLoginWithoutTTY(c *check.C) {
	cmd := exec.Command(dockerBinary, "login")

	// Send to stdin so the process does not get the TTY
	cmd.Stdin = bytes.NewBufferString("buffer test string \n")

	// run the command and block until it's done
	err := cmd.Run()
	c.Assert(err, checker.NotNil) //"Expected non nil err when loginning in & TTY not available"
}

func (s *DockerRegistriesSuite) TestLoginAgainstDefaultRegistries(c *check.C) {
	c.Assert(s.d.StartWithBusybox("--add-registry="+s.regWithAuth.url), check.IsNil)

	// check we can login against the first default registry which isn't docker.io
	out, err := s.d.Cmd("login", "-u", s.regWithAuth.username, "-p", s.regWithAuth.password, "-e", s.regWithAuth.email)
	c.Assert(err, check.IsNil, check.Commentf(out))

	// check we can login against "docker.io" which translates to their v1 server address
	// check we can login against the first default registry which isn't docker.io
	wanted := "Wrong login/password, please try again"
	out, err = s.d.Cmd("login", "-u", s.regWithAuth.username, "-p", s.regWithAuth.password, "-e", s.regWithAuth.email, registry.IndexName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	if !strings.Contains(out, wanted) {
		c.Fatalf("wanted %s, got %s", wanted, out)
	}

	// stop the daemon
	// restart it w/o --add-registry and check I can attempt login against "docker.io"
	c.Assert(s.d.Stop(), check.IsNil)
	c.Assert(s.d.StartWithBusybox(), check.IsNil)

	out, err = s.d.Cmd("login", "-u", s.regWithAuth.username, "-p", s.regWithAuth.password, "-e", s.regWithAuth.email, registry.IndexName)
	c.Assert(err, check.NotNil, check.Commentf(out))
	if !strings.Contains(out, wanted) {
		c.Fatalf("wanted %s, got %s", wanted, out)
	}
}
