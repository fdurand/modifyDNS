package scutil

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/davecgh/go-spew/spew"
	utilexec "k8s.io/utils/exec"
)

const (
	cmdScutil       string = "scutil"
	cmdNetworksetup string = "networksetup"
)

// Interface is an injectable interface for running netsh commands.  Implementations must be goroutine-safe.
type Interface interface {
	// GetDNSServers retreive the dns servers
	// GetDNSServers(args []string) (bool, error)
	GetDNSServers(iface string) DNSConfig
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
func (runner *runner) GetDNSServers(ifname string) DNSConfig {
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
	runner.InterfaceAliasName(currentInterface.IfIndex)

	return currentInterface
}

func (runner *runner) InterfaceAliasName(iface string) (string, error) {

	args := []string{
		"-listnetworkserviceorder",
	}

	output, _ := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	spew.Dump(outputLines)

	interfacePattern := regexp.MustCompile("\\(Hardware Port:\\s+(.*),\\s+Device:\\s+(.*)\\)")

	for _, outputLine := range outputLines {
		if strings.Contains(outputLine, "Hardware Port") {
			match := interfacePattern.FindStringSubmatch(outputLine)
			spew.Dump(match)
		} else {
			continue
		}
	}
	return iface, nil
}

// Set DNS server on the interface (name or index)
func (runner *runner) SetDNSServer(iface string, dns string) error {
	args := []string{
		"interface", "ipv4", "set", "dnsservers", "name=" + strconv.Quote(iface), "source=static", strconv.Quote(dns), "primary",
	}
	cmd := strings.Join(args, " ")
	if stdout, err := runner.exec.Command(cmdScutil, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set dns servers on [%v], error: %v. cmd: %v. stdout: %v", iface, err.Error(), cmd, string(stdout))
	}

	return nil
}

// Reset DNS on the interface (name or index)
func (runner *runner) ResetDNSServer(iface string) error {
	args := []string{
		"interface", "ipv4", "set", "dnsservers", "name=" + strconv.Quote(iface), "source=dhcp",
	}
	cmd := strings.Join(args, " ")
	if stdout, err := runner.exec.Command(cmdScutil, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reset dns servers on [%v], error: %v. cmd: %v. stdout: %v", iface, err.Error(), cmd, string(stdout))
	}

	return nil
}
