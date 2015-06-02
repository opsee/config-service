package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"sort"

	"github.com/coreos/go-etcd/etcd"
)

const (
	configRoot = "/opsee.co"
	cliName    = "cf"
)

var (
	// It's a short goddamn list. Just sort it yourself.
	commands = []string{
		"get",
		"set",
		"current",
	}

	client *etcd.Client
)

// Make this look more pretty like
// https://github.com/coreos/fleet/blob/master/fleetctl/fleetctl.go

func stderr(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func stdout(format string, args ...interface{}) {
	fmt.Fprintf(os.Stdout, format, args...)
}

func etcdKey(service string, tag string) string {
	return fmt.Sprintf("%s/%s/%s", configRoot, service, tag)
}

func getCurrentVersion(service string) (string, error) {
	response, err := client.Get(etcdKey(service, "current"), false, false)
	if err != nil {
		return "", err
	}

	return response.Node.Value, nil
}

func getConfig(service string, tag string) (string, error) {
	response, err := client.Get(etcdKey(service, tag), false, false)
	if err != nil {
		return "", err
	}

	return response.Node.Value, nil
}

func setConfig(service string, tag string) (string, error) {
	configReader := bufio.NewReader(os.Stdin)

	configBytes, err := ioutil.ReadAll(configReader)

	if err != nil {
		return "", err
	}

	configStr := string(configBytes)

	response, err := client.Set(etcdKey(service, tag), configStr, 0)
	if err != nil {
		return "", err
	}

	responseJSON, err := json.Marshal(response)
	if err != nil {
		return "", err
	}

	fmt.Println(responseJSON)
	return configStr, nil
}

func setCurrent(service string, tag string) (string, error) {
	if tag == "" {
		stderr("%v: must supply tag to set as current version", cliName)
		os.Exit(2)
	}

	response, err := client.Set(etcdKey(service, "current"), tag, 0)
	if err != nil {
		return "", err
	}

	return response.Node.Value, nil
}

func stringInArray(a []string, c string) bool {
	idx := sort.SearchStrings(a, c) // Must be sorted
	return idx < len(a)
}

func getServices(c *etcd.Client) []string {
	response, err := c.Get(configRoot, false, false)
	if err != nil {
		log.Fatal(err)
	}

	nodes := response.Node.Nodes
	nodeCount := nodes.Len()

	var services = make([]string, nodeCount)

	for i := 0; i < nodeCount; i++ {
		services[i] = nodes[i].Value
	}

	return services
}

func usage() {
	helpfmt := `usage: %v [command] [service] [tag]

Commands
get				get a version of a configuration file
set				set a version of a configuration file
current   set the current version of a configuration file

Get
Gets the configuration for the specified service. If unspecified,
it returns the current configuration.

Set
Sets the configuration for the specified version from STDIN. You
MUST specify a tag.

Current
Sets the current version to the specified tag for the given service.
`
	fmt.Printf(helpfmt, cliName)
}

// I am so fucking lazy sometimes.
func main() {
	machines := []string{"http://127.0.0.1:2379"}
	client = etcd.NewClient(machines)

	args := os.Args[1:]

	if len(args) < 2 {
		usage()
		os.Exit(2)
	}

	var tag string

	command := args[0]
	service := args[1]
	if len(args) == 3 {
		tag = args[2]
	}

	if command == "" || !stringInArray(commands, command) {
		stderr("%v: unkown command: %s", cliName, command)
		os.Exit(2)
	}

	if service == "" {
		stderr("%v: must specify a service", cliName)
		os.Exit(2)
	}

	// Command dispatch should be better. :(
	switch command {
	case "get":
		if len(args) == 2 {
			current, err := getCurrentVersion(service)
			tag = current
			if err != nil {
				stderr("%v: %v", cliName, err)
				os.Exit(2)
			}
		}
		if config, err := getConfig(service, tag); err != nil {
			stderr("%v: %v", cliName, err)
			os.Exit(2)
		} else {
			fmt.Println(config)
		}
	case "set":
		if len(args) != 3 {
			stderr("%v: must supply tag to set", cliName)
			os.Exit(2)
		}
		if _, err := setConfig(service, tag); err != nil {
			stderr("%v: %v", cliName, err)
			os.Exit(2)
		}
	case "current":
		if _, err := setCurrent(service, tag); err != nil {
			stderr("%v: %v", cliName, err)
			os.Exit(2)
		}
	}
}
