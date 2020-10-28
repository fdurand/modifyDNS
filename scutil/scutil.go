package scutil

import (
	"regexp"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	utilexec "k8s.io/utils/exec"
)

const (
	cmdScutil string = "scutil"
)

// Interface is an injectable interface for running netsh commands.  Implementations must be goroutine-safe.
type Interface interface {
	// GetDNSServers retreive the dns servers
	// GetDNSServers(args []string) (bool, error)
	GetDNSServers(iface string)
	// Set DNS server on this interface (name or index)
	// SetDNSServer(iface string, dns string) error
	// // Reset DNS server on this interface (name or index)
	// ResetDNSServer(iface string) error
}

// runner implements Interface in terms of exec("netsh").
type runner struct {
	mu   sync.Mutex
	exec utilexec.Interface
}

type DNSConfig struct {
	Domain       string
	SearchDomain []string
	NameServers  []string
	IfIndex      string
	Flags        string
	Reach        string
	Options      string
}

// New returns a new Interface which will exec scutil.
func New(exec utilexec.Interface) Interface {

	if exec == nil {
		exec = utilexec.New()
	}

	runner := &runner{
		exec: exec,
	}

	return runner
}

// GetDNSServers uses the show addresses command and returns a formatted structure
func (runner *runner) GetDNSServers(ifname string) {
	args := []string{
		"--dns",
	}

	output, _ := runner.exec.Command(cmdScutil, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	interfacePattern := regexp.MustCompile("^\\d+\\s+\\((.*)\\)")

	currentInterface := DNSConfig{}

	found := false
	for _, outputLine := range outputLines {
		if !found {
			spew.Dump(outputLine)
			if strings.Contains(outputLine, "DNS configuration (for scoped queries)") {
				found = true
			} else {
				continue
			}
		}

		parts := strings.SplitN(outputLine, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if strings.HasPrefix(key, "if_index") {
			match := interfacePattern.FindStringSubmatch(value)
			if match[1] == ifname {
				found = true
				currentInterface.IfIndex = ifname
			}
		} else if strings.HasPrefix(key, "nameserver") {
			currentInterface.NameServers = append(currentInterface.NameServers, value)
		} else if strings.HasPrefix(key, "search domain") {
			currentInterface.SearchDomain = append(currentInterface.SearchDomain, value)
		} else if strings.HasPrefix(key, "flags") {
			currentInterface.Flags = value
		} else if strings.HasPrefix(key, "reach") {
			currentInterface.Reach = value
		} else if strings.HasPrefix(key, "domain") {
			currentInterface.Domain = value
		} else if strings.HasPrefix(key, "reach") {
			currentInterface.Reach = value
		} else if strings.HasPrefix(key, "options") {
			currentInterface.Options = value
		}
	}
	spew.Dump(currentInterface)
}
