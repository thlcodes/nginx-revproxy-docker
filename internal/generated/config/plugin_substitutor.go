// Code generated by genfig plugin 'substitutor'; DO NOT EDIT.

package config

import (
	"strings"
)

var _ = strings.Contains

const (
	maxSubstitutionIteraions = 5
)

var (
	raw Config
)

// Substitute replaces all.
// The return value informs, whether all substitutions could be
// applied within {maxSubstitutionIteraions} or not
func (c *Config) Substitute() bool {
	c.ResetSubstitution()

	// backup the "raw" configuration
	raw = *c

	run := 0
	for {
		if run == maxSubstitutionIteraions {
			return false
		}
		if c.substitute() == 0 {
			return true
		}
		run += 1
	}
}

// ResetSubstitution resets the configuration to the state,
// before the substitution was applied
func (c *Config) ResetSubstitution() {
	c = &raw
}

// substitute tries to replace all substitutions in strings
func (c *Config) substitute() int {
	cnt := 0

	r := strings.NewReplacer(
		"${apipath}", c.Apipath,

		"${basepath}", c.Basepath,

		"${greeter.defaultname}", c.Greeter.DefaultName,

		"${greeter.hello}", c.Greeter.Hello,

		"${host}", c.Host,

		"${registry}", c.Registry,

		"${service.name}", c.Service.Name,

		"${service.version}", c.Service.Version,

		"${services.storage.addr}", c.Services.Storage.Addr,

		"${tracing.secret}", c.Tracing.Secret,

		"${tracing.service}", c.Tracing.Service,

		"${tracing.uri}", c.Tracing.Uri,
	)

	if strings.Contains(c.Apipath, "${") {
		cnt += 1
		c.Apipath = r.Replace(c.Apipath)
		if !strings.Contains(c.Apipath, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Basepath, "${") {
		cnt += 1
		c.Basepath = r.Replace(c.Basepath)
		if !strings.Contains(c.Basepath, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Greeter.DefaultName, "${") {
		cnt += 1
		c.Greeter.DefaultName = r.Replace(c.Greeter.DefaultName)
		if !strings.Contains(c.Greeter.DefaultName, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Greeter.Hello, "${") {
		cnt += 1
		c.Greeter.Hello = r.Replace(c.Greeter.Hello)
		if !strings.Contains(c.Greeter.Hello, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Host, "${") {
		cnt += 1
		c.Host = r.Replace(c.Host)
		if !strings.Contains(c.Host, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Registry, "${") {
		cnt += 1
		c.Registry = r.Replace(c.Registry)
		if !strings.Contains(c.Registry, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Service.Name, "${") {
		cnt += 1
		c.Service.Name = r.Replace(c.Service.Name)
		if !strings.Contains(c.Service.Name, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Service.Version, "${") {
		cnt += 1
		c.Service.Version = r.Replace(c.Service.Version)
		if !strings.Contains(c.Service.Version, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Services.Storage.Addr, "${") {
		cnt += 1
		c.Services.Storage.Addr = r.Replace(c.Services.Storage.Addr)
		if !strings.Contains(c.Services.Storage.Addr, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Tracing.Secret, "${") {
		cnt += 1
		c.Tracing.Secret = r.Replace(c.Tracing.Secret)
		if !strings.Contains(c.Tracing.Secret, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Tracing.Service, "${") {
		cnt += 1
		c.Tracing.Service = r.Replace(c.Tracing.Service)
		if !strings.Contains(c.Tracing.Service, "${") {
			cnt -= 1
		}
	}

	if strings.Contains(c.Tracing.Uri, "${") {
		cnt += 1
		c.Tracing.Uri = r.Replace(c.Tracing.Uri)
		if !strings.Contains(c.Tracing.Uri, "${") {
			cnt -= 1
		}
	}

	return cnt
}