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
	GetDNSServers()
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
func (runner *runner) GetDNSServers() {
	args := []string{
		"--dns",
	}

	output, err := runner.exec.Command(cmdScutil, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	interfacePattern := regexp.MustCompile("^\\d+\\s+\\((.*)\\)")

	for _, outputLine := range outputLines {
		parts := strings.SplitN(outputLine, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if strings.HasPrefix(key, "if_index") {
			match := interfacePattern.FindStringSubmatch(value)
			spew.Dump(match[1])
		}
	}
	spew.Dump(err)
}
