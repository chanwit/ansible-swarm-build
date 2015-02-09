package main

import "fmt"

type machine struct {
}

type machines []machine

func (m machines) Install(what string, commit string) {
	fmt.Printf("%s %s\n", what, commit)
}

func (m machines) Run(cmd string, args ...string) {
	fmt.Printf("%d: %s %s\n", len(m), cmd, args)
}

func (m machine) Run(cmd string, args ...string) {
	fmt.Printf("%s %s\n", cmd, args)
}

type instance int
type mb int

var Instances = instance(1)
var MB = mb(1)

func NewCluster(num instance, image string, mem mb) machines {
	fmt.Printf("%d %s %d\n", num, image, mem)
	m := make([]machine, int(num))
	for i := 1; i <= int(num); i++ {
		m[i-1] = machine{}
        } 
	return machines(m)
}

func main() {
	m := NewCluster(4*Instances, "ubuntu", 512*MB)
        // m.InstallDocker("1.4.1")
	// m.InstallSwarm("github.com/docker/swarm", "0a0b0d0e")
        machines(m[1:4]).Run("swarm", "join", "etcd://master")
        m[0].Run("swarm", "manage")	

}
