// Package shell generate shell script
package shell

import (
	"fmt"
	"strings"

	l "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/tinynetwork/tinet/internal/pkg/utils"
)

var log = l.New()

// Tn tinet config
type Tn struct {
	PreCmd      []PreCmd     `yaml:"precmd"`
	PreInit     []PreInit    `yaml:"preinit"`
	PostInit    []PostInit   `yaml:"postinit"`
	PostFini    []PostFini   `yaml:"postfini"`
	Nodes       []Node       `yaml:"nodes" mapstructure:"nodes"`
	Switches    []Switch     `yaml:"switches" mapstructure:"switches"`
	NodeConfigs []NodeConfig `yaml:"node_configs" mapstructure:"node_configs"`
	Test        []Test       `yaml:"test"`
}

// PreCmd
type PreCmd struct {
	// Cmds []Cmd `yaml:"cmds"`
	Cmds []Cmd `yaml:"cmds" mapstructure:"cmds"`
}

// PreInit
type PreInit struct {
	Cmds []Cmd `yaml:"cmds" mapstructure:"cmds"`
}

// PostInit
type PostInit struct {
	Cmds []Cmd `yaml:"cmds" mapstructure:"cmds"`
}

// PostFini
type PostFini struct {
	Cmds []Cmd `yaml:"cmds" mapstructure:"cmds"`
}

// Node
type Node struct {
	Name       string      `yaml:"name" mapstructure:"name"`
	Type       string      `yaml:"type" mapstructure:"type"`
	NetBase    string      `yaml:"net_base" mapstructure:"net_base"`
	Image      string      `yaml:"image" mapstructure:"image"`
	Interfaces []Interface `yaml:"interfaces" mapstructure:"interfaces"`
	Sysctls    []Sysctl    `yaml:"sysctls" mapstructure:"sysctls"`
	Mounts     []string    `yaml:"mounts,flow" mapstructure:"mounts,flow"`
}

// Interface
type Interface struct {
	Name string `yaml:"name"`
	Type string `yaml:"type"`
	Args string `yaml:"args"`
	Addr string `yaml:"addr"`
}

// Sysctl
type Sysctl struct {
	Sysctl string `yaml:"string"`
}

// Switch
type Switch struct {
	Name       string      `yaml:"name"`
	Interfaces []Interface `yaml:"interfaces" mapstructure:"interfaces"`
}

// NodeConfig
type NodeConfig struct {
	Name string `yaml:"name"`
	Cmds []Cmd  `yaml:"cmds" mapstructure:"cmds"`
}

// Cmd
type Cmd struct {
	Cmd string `yaml:"cmd"`
}

// Test
type Test struct {
	Name string
	Cmds []Cmd `yaml:"cmds" mapstructure:"cmds"`
}

// BuildCmd
func BuildCmd(nodes []Node) string {
	return "sorry not implement..."
}

// ExecConf Execute NodeConfig command
func (nodeConfig *NodeConfig) ExecConf(nodeType string) []string {
	var execConfCmds []string
	for _, nodeConfigCmd := range nodeConfig.Cmds {
		var execConfCmd string
		if nodeType == "docker" {
			execConfCmd = fmt.Sprintf("docker exec %s %s", nodeConfig.Name, nodeConfigCmd.Cmd)
		} else if nodeType == "netns" {
			execConfCmd = fmt.Sprintf("ip netns exec %s %s", nodeConfig.Name, nodeConfigCmd.Cmd)
		} else if nodeType == "" {
			execConfCmd = fmt.Sprintf("docker exec %s %s", nodeConfig.Name, nodeConfigCmd.Cmd)
		} else {
			// err := fmt.Errorf("not supported node type...")
			// log.Fatal(err)
			return []string{""}
		}
		// fmt.Println(execConfigCmd)
		execConfCmds = append(execConfCmds, execConfCmd)
	}

	return execConfCmds
}

// DeleteNode Delete docker and netns
func (node *Node) DeleteNode() []string {
	var deleteCmd string
	if node.Type == "docker" || node.Type == "" {
		deleteCmd = fmt.Sprintf("docker rm -f %s", node.Name)
	} else if node.Type == "netns" {
		deleteCmd = fmt.Sprintf("ip netns del %s", node.Name)
	} else {
		// err := fmt.Errorf("not supported node type...")
		// log.Fatal(err)
		return []string{""}
	}

	deleteNsCmd := fmt.Sprintf("rm -rf /var/run/netns/%s", node.Name)

	deleteNodeCmds := []string{deleteCmd, deleteNsCmd}

	return deleteNodeCmds
}

// DeleteSwitch Delete bridge
func (br *Switch) DeleteSwitch() string {
	deleteBrCmd := fmt.Sprintf("ip link delete %s", br.Name)
	return deleteBrCmd
}

// Exec Select Node exec command
func (tnconfig *Tn) Exec(nodeName string, Cmds []string) string {
	var execCommand string
	var selectedNode *Node

	for i, node := range tnconfig.Nodes {
		if node.Name == nodeName {
			selectedNode = &tnconfig.Nodes[i]
		}
	}

	if selectedNode == nil {
		err := fmt.Errorf("no such node...\n")
		log.Error(err)
	}

	if selectedNode.Type == "docker" {
		execCommand = fmt.Sprintf("docker exec %s", nodeName)
	} else if selectedNode.Type == "netns" {
		execCommand = fmt.Sprintf("ip netns exec %s", nodeName)
	} else if selectedNode.Type == "" {
		execCommand = fmt.Sprintf("docker exec %s", nodeName)
	} else {
		err := fmt.Errorf("no such node type...\n")
		log.Error(err)
	}

	var cmdStr string
	for _, cmd := range Cmds {
		cmdStr += fmt.Sprintf(" %s", cmd)
		execCommand += cmdStr
	}

	return execCommand
}

// GenerateFile Generate tinet template config file
func GenerateFile() (string, error) {
	precmd := PreCmd{
		Cmds: []Cmd{
			Cmd{
				Cmd: "",
			},
		},
	}

	preinit := PreInit{
		Cmds: []Cmd{
			Cmd{
				Cmd: "",
			},
		},
	}
	postinit := PostInit{
		Cmds: []Cmd{
			Cmd{
				Cmd: "echo hoge",
			},
		},
	}

	postfini := PostFini{
		Cmds: []Cmd{
			Cmd{
				Cmd: "",
			},
		},
	}

	nodes := Node{
		Name:    "R1",
		Image:   "slankdev/ubuntu18.04",
		NetBase: "",
		Interfaces: []Interface{
			Interface{
				Name: "net0",
				Type: "direct",
				Args: "C1#net0",
			},
			Interface{
				Name: "net1",
				Type: "bridge",
				Args: "B0",
			},
			Interface{
				Name: "net2",
				Type: "veth",
				Args: "peer0",
			},
			Interface{
				Name: "net3",
				Type: "phys",
			},
		},
	}

	switches := Switch{
		Name: "B0",
		Interfaces: []Interface{
			Interface{
				Name: "net0",
				Type: "docker",
				Args: "R1",
			},
			Interface{
				Name: "net0",
				Type: "netns",
				Args: "R2",
			},
		},
	}

	nodeConfigs := []NodeConfig{
		NodeConfig{
			Name: "C0",
			Cmds: []Cmd{
				Cmd{
					Cmd: "ip link set dev net0 up",
				},
			},
		},
		NodeConfig{
			Name: "C1",
			Cmds: []Cmd{
				Cmd{
					Cmd: "echo slankdev slankdev",
				},
				Cmd{
					Cmd: "echo slankdev &&\necho slankdev",
				},
			},
		},
	}

	tests := []Test{
		Test{
			Name: "p2p",
			Cmds: []Cmd{
				Cmd{
					Cmd: "docker exec C0 ping -c2 10.0.0.2",
				},
				Cmd{
					Cmd: "echo hoge",
				},
			},
		},
		Test{
			Name: "lo",
			Cmds: []Cmd{
				Cmd{
					Cmd: "docker exec C0 ping -c2 10.255.0.1",
				},
				Cmd{
					Cmd: "echo hoge",
				},
			},
		},
	}

	tnconfig := &Tn{
		PreCmd:      []PreCmd{precmd},
		PreInit:     []PreInit{preinit},
		PostInit:    []PostInit{postinit},
		PostFini:    []PostFini{postfini},
		Nodes:       []Node{nodes},
		Switches:    []Switch{switches},
		NodeConfigs: nodeConfigs,
		Test:        tests,
	}

	data, err := yaml.Marshal(tnconfig)
	if err != nil {
		return "", err
	}
	//
	// err = ioutil.WriteFile(cfgFile, data, 0644)
	// if err != nil {
	// 	return err
	// }

	// return nil
	return string(data), nil
}

// DockerPs Show docker ps
func DockerPs(all bool) string {
	dockerPsCmd := "docker ps --format 'table {{.ID}}\\t{{.Names}}\\t{{.Status}}\\t{{.Networks}}'"
	if all {
		dockerPsCmd += " -a"
	}

	return dockerPsCmd
}

// NetnsPs Show ip netns list
func NetnsPs() string {
	netnsListCmd := "ip netns list"

	return netnsListCmd
}

// Pull pull Docker Image
func Pull(nodes []Node) []string {
	var images []string
	for _, node := range nodes {
		images = append(images, node.Image)
	}
	images = utils.RemoveDuplicatesString(images)

	var pullCmds []string
	for _, image := range images {
		pullCmd := fmt.Sprintf("docker pull %s", image)
		pullCmds = append(pullCmds, pullCmd)
	}

	return pullCmds
}

// TnTestCmdExec Execute test cmds
func TnTestCmdExec(tests []Test) []string {
	var tnTestCmds []string
	for _, test := range tests {

		for _, testCmd := range test.Cmds {
			tnTestCmds = append(tnTestCmds, testCmd.Cmd)
		}

	}

	return tnTestCmds
}

// ExecCmd Execute cmds
func ExecCmd(cmds []Cmd) []string {
	var execCmds []string
	for _, cmd := range cmds {
		execCmds = append(execCmds, cmd.Cmd)
	}

	return execCmds
}

// CreateNode Create nodes set in config
func (node *Node) CreateNode() []string {
	var createNodeCmds []string

	var createNodeCmd string
	if node.NetBase == "" {
		node.NetBase = "none"
	}
	if node.Type == "docker" || node.Type == "" {
		createNodeCmd = fmt.Sprintf("docker run -td --hostname %s --net %s --name %s --rm --privileged ", node.Name, node.NetBase, node.Name)
		if len(node.Sysctls) != 0 {
			for _, sysctl := range node.Sysctls {
				createNodeCmd += fmt.Sprintf("--sysctl %s ", sysctl.Sysctl)
			}
		}

		if len(node.Mounts) != 0 {
			for _, mount := range node.Mounts {
				createNodeCmd += fmt.Sprintf("-v %s ", mount)
			}
		}

		createNodeCmd += node.Image
	} else if node.Type == "netns" {
		createNodeCmd = fmt.Sprintf("ip netns add %s", node.Name)
	} else {
		// err := fmt.Errorf("unknown nodetype %s", node.Type)
		// log.Fatal(err)
		createNodeCmd = fmt.Sprintf("unknown nodetype %s", node.Type)
	}

	createNodeCmds = append(createNodeCmds, createNodeCmd)

	if node.Type == "netns" {
		if len(node.Sysctls) != 0 {
			for _, sysctl := range node.Sysctls {
				sysctlNsCmd := fmt.Sprintf("ip netns exec %s sysctl -w %s", node.Name, sysctl.Sysctl)
				createNodeCmds = append(createNodeCmds, sysctlNsCmd)
			}
		}
		infloUpCmd := fmt.Sprintf("ip netns exec %s ip link set lo up", node.Name)
		createNodeCmds = append(createNodeCmds, infloUpCmd)
	}

	return createNodeCmds
}

// HostLinkUp Link up link of host
func HostLinkUp(linkName string) string {
	linkUpCmd := fmt.Sprintf("ip link set %s up", linkName)
	return linkUpCmd
}

// NetnsLinkUp Link up link of netns
func NetnsLinkUp(netnsName string, linkName string) string {
	netnsLinkUpCmd := fmt.Sprintf("ip netns exec %s ip link set %s up", netnsName, linkName)
	return netnsLinkUpCmd
}

// CreateSwitch Create bridge set in config
func (bridge *Switch) CreateSwitch() []string {
	var createSwitchCmds []string

	addSwitchCmd := fmt.Sprintf("ip link add %s type bridge", bridge.Name)
	createSwitchCmds = append(createSwitchCmds, addSwitchCmd)

	bridgeUpCmd := HostLinkUp(bridge.Name)
	createSwitchCmds = append(createSwitchCmds, bridgeUpCmd)

	return createSwitchCmds
}

// N2nLink Connect links between nodes
func (inf *Interface) N2nLink(nodeName string) []string {
	var n2nLinkCmds []string

	nodeinf := inf.Name
	peerinfo := strings.Split(inf.Args, "#")
	peerNode := peerinfo[0]
	peerinf := peerinfo[1]
	n2nlinkCmd := fmt.Sprintf("ip link add %s netns %s type veth peer name %s netns %s", nodeinf, nodeName, peerinf, peerNode)
	n2nLinkCmds = append(n2nLinkCmds, n2nlinkCmd)
	n2nLinkCmds = append(n2nLinkCmds, NetnsLinkUp(nodeName, nodeinf))
	n2nLinkCmds = append(n2nLinkCmds, NetnsLinkUp(peerNode, peerinf))

	if len(inf.Addr) != 0 {
		addrSetCmd := fmt.Sprintf("ip netns exec %s ip link set %s address %s", nodeName, inf.Name, inf.Addr)
		n2nLinkCmds = append(n2nLinkCmds, addrSetCmd)
	}

	return n2nLinkCmds
}

// S2nLink Connect links between nodes and switches
func (inf *Interface) S2nLink(nodeName string) []string {
	var s2nLinkCmds []string

	nodeinf := inf.Name
	peerBr := inf.Args
	peerBrInf := fmt.Sprintf("%s-%s", peerBr, nodeName)
	s2nLinkCmd := fmt.Sprintf("ip link add %s netns %s type veth peer name %s", nodeinf, nodeName, peerBrInf)
	s2nLinkCmds = append(s2nLinkCmds, s2nLinkCmd)
	s2nLinkCmds = append(s2nLinkCmds, NetnsLinkUp(nodeName, nodeinf))
	s2nLinkCmds = append(s2nLinkCmds, HostLinkUp(peerBrInf))
	setBrLinkCmd := fmt.Sprintf("ip link set dev %s master %s", peerBrInf, peerBr)
	s2nLinkCmds = append(s2nLinkCmds, setBrLinkCmd)

	return s2nLinkCmds
}

// V2cLink Connect links between veth and container
func (inf *Interface) V2cLink(nodeName string) []string {
	var v2cLinkCmds []string
	nodeinf := inf.Name
	peerName := inf.Args
	v2cLinkCmd := fmt.Sprintf("ip link add %s type veth peer name %s", nodeinf, peerName)
	v2cLinkCmds = append(v2cLinkCmds, v2cLinkCmd)

	v2cLinkCmds = append(v2cLinkCmds, inf.P2cLink(nodeName)...)
	v2cLinkCmds = append(v2cLinkCmds, HostLinkUp(peerName))

	return v2cLinkCmds
}

// P2cLink Connect links between phys-eth and container
func (inf *Interface) P2cLink(nodeName string) []string {
	var p2cLinkCmds []string

	physInf := inf.Name
	setNsCmd := fmt.Sprintf("ip link set dev %s netns %s", physInf, nodeName)
	p2cLinkCmds = append(p2cLinkCmds, setNsCmd)
	execNsCmd := fmt.Sprintf("ip netns exec %s ip link set %s up", nodeName, physInf)
	p2cLinkCmds = append(p2cLinkCmds, execNsCmd)
	delNsCmd := fmt.Sprintf("ip netns del %s", nodeName)
	p2cLinkCmds = append(p2cLinkCmds, delNsCmd)

	return p2cLinkCmds
}

// Mount_docker_netns Mount docker netns to ip netns
func (node *Node) Mount_docker_netns() []string {
	var mountDockerNetnsCmds []string

	netnsDir := "/var/run/netns"
	mkdirCmd := fmt.Sprintf("mkdir -p %s", netnsDir)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, mkdirCmd)
	dockerPid := GetContainerPid(node.Name)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, dockerPid)
	mountDockerNetnsCmd := fmt.Sprintf("ln -s /proc/$PID/ns/net /var/run/netns/%s", node.Name)
	mountDockerNetnsCmds = append(mountDockerNetnsCmds, mountDockerNetnsCmd)

	return mountDockerNetnsCmds
}

// GetContainerPid func is Output get Docker PID Command
func GetContainerPid(nodename string) string {
	getpidcmd := fmt.Sprintf("PID=`docker inspect %s --format '{{.State.Pid}}'`", nodename)

	return getpidcmd
}
