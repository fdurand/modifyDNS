package scutil

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"

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
	GetDNSServers(iface string) error
	// Set DNS server
	SetDNSServer(dns string) error
	// // Reset DNS server
	ResetDNSServer() error
}

// runner implements Interface in terms of exec("netsh").
type runner struct {
	mu                 sync.Mutex
	exec               utilexec.Interface
	InterFaceDNSConfig DNSConfig
}

type DNSConfig struct {
	Domain       string
	SearchDomain []string
	NameServers  []string
	IfIndex      string
	IfName       string
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
func (runner *runner) GetDNSServers(ifname string) error {
	args := []string{
		"--dns",
	}

	output, _ := runner.exec.Command(cmdScutil, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	interfacePattern := regexp.MustCompile("^\\d+\\s+\\((.*)\\)")

	runner.InterFaceDNSConfig = DNSConfig{}

	// currentInterface := DNSConfig{}

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
				runner.InterFaceDNSConfig.IfIndex = ifname
			}
		} else if strings.HasPrefix(key, "nameserver") {
			runner.InterFaceDNSConfig.NameServers = append(runner.InterFaceDNSConfig.NameServers, value)
		} else if strings.HasPrefix(key, "search domain") {
			runner.InterFaceDNSConfig.SearchDomain = append(runner.InterFaceDNSConfig.SearchDomain, value)
		} else if strings.HasPrefix(key, "flags") {
			runner.InterFaceDNSConfig.Flags = value
		} else if strings.HasPrefix(key, "reach") {
			runner.InterFaceDNSConfig.Reach = value
		} else if strings.HasPrefix(key, "domain") {
			runner.InterFaceDNSConfig.Domain = value
		} else if strings.HasPrefix(key, "reach") {
			runner.InterFaceDNSConfig.Reach = value
		} else if strings.HasPrefix(key, "options") {
			runner.InterFaceDNSConfig.Options = value
		}
	}

	err := runner.InterfaceAliasName()

	return err
}

func (runner *runner) InterfaceAliasName() error {

	args := []string{
		"-listnetworkserviceorder",
	}

	output, _ := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput()

	DNSString := string(output[:])

	outputLines := strings.Split(DNSString, "\n")

	interfacePattern := regexp.MustCompile("\\(Hardware Port:\\s+(.*),\\s+Device:\\s+(.*)\\)")

	err := errors.New("Unable to find the interface alias")

	for _, outputLine := range outputLines {
		if strings.Contains(outputLine, "Hardware Port") {
			match := interfacePattern.FindStringSubmatch(outputLine)
			if match[2] == runner.InterFaceDNSConfig.IfIndex {
				runner.InterFaceDNSConfig.IfName = match[1]
				err = nil
			}
		} else {
			continue
		}
	}
	return err
}

// Set DNS server on the interface (name or index)
func (runner *runner) SetDNSServer(dns string) error {
	args := []string{
		"-setdnsservers", runner.InterFaceDNSConfig.IfName, dns,
	}
	cmd := strings.Join(args, " ")
	if stdout, err := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to set dns servers on [%v], error: %v. cmd: %v. stdout: %v", runner.InterFaceDNSConfig.IfName, err.Error(), cmd, string(stdout))
	}
	return nil
}

// Reset DNS on the interface (name or index)
func (runner *runner) ResetDNSServer() error {
	args := []string{
		"-setdnsservers", runner.InterFaceDNSConfig.IfName, strings.Join(runner.InterFaceDNSConfig.NameServers[:], " "),
	}
	cmd := strings.Join(args, " ")

	if stdout, err := runner.exec.Command(cmdNetworksetup, args...).CombinedOutput(); err != nil {
		return fmt.Errorf("failed to reset dns servers on [%v], error: %v. cmd: %v. stdout: %v", runner.InterFaceDNSConfig.IfName, err.Error(), cmd, string(stdout))
	}

	return nil
}
