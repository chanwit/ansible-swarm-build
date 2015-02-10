package main

import "fmt"
import "os/exec"

type instance int
type mb int

var Instances = instance(1)
var MB = mb(1)

type machine struct {
	name string
	uid  string
}

type machines []machine

func (m machines) Install(what string, commit string) {
	fmt.Printf("len(m): %d\n", len(m))
	for _, mm := range m {
		mm.Install(what, commit)
	}
}

func (m machine) Install(what string, commit string) {
	fmt.Printf("%s %s\n", what, commit)
}

func (m machines) Run(cmd string, args ...string) {
	fmt.Printf("len(m): %d\n", len(m))
	for _, mm := range m {
		mm.Run(cmd, args...)
	}
}

func (m machine) Run(cmd string, args ...string) {
	fmt.Printf("%s %s\n", cmd, args)
}

func Import() {
	name := fmt.Sprintf("%s%00d", prefix, i)
	cmd := exec.Command("vboxmanage", "import", file,
		"--vsys", "0", "--vmname", name,
		"--memory", fmt.Sprintf("%d", memory),
		"--unit", "10", "--disk", fmt.Sprintf("%s/%s.vmdk", pwd, name))
	if err := cmd.Run(); err != nil {
		return err
	}
	out, err := cmd.Output()
	// parse
}

func NewCluster(num instance, image string, mem mb) machines {
	// import 1
	// set 2 interface
	// 1. host-only (for talking to master)
	// 2. nat (for internet connection)
	// clone n-1
	fmt.Printf("%d %s %d\n", num, image, mem)
	m := make([]machine, int(num))
	for i := 1; i <= int(num); i++ {
		m[i-1] = machine{}
	}
	return machines(m)
}

func main() {
	m := NewCluster(4*Instances,
		"ubuntu",
		512*MB,
		map[string]string{
			"SWARM_DISCOVERY": "etcd://master",
			"SWARM_HOST": "0.0.0.0:2375",
		})
	m.Install("docker", "1.4.1")
	m.Install("github.com/docker/swarm", "0a0b0d0e")

	master, slaves := m[0], machines(m[1:])

	master.Install("etcd", "1.0")
	master.Install("zookeeper", "3.4.3")

	slaves.Run("swarm", "join", "--addr", "{{eth0.ipv4.address}}:2375", discovery)
	master.Run("swarm", "manage", discovery)

	master.Run("docker", "run", "nginx")
}
