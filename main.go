package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var VBOX_MANAGE = "VBoxManage"

func Import(file string, memory int) error {
	name := "base"
	cmd := exec.Command(VBOX_MANAGE, "import", os.Getenv("PWD")+"/"+file,
		"--vsys", "0", "--vmname", name,
		"--memory", fmt.Sprintf("%d", memory),
	)
	_, err := cmd.Output()
	if err == nil {
		log.Infof("Imported \"%s\"", name)
	}
	if err != nil {
		return err
	}

	_, err = exec.Command(VBOX_MANAGE, "snapshot", "base", "take", "origin").Output()
	log.Info("Snapshot \"origin\" taken")
	return err
}

func Modify(name string, adapter string, i int) error {
	err := exec.Command(VBOX_MANAGE, "modifyvm", name,
		"--natpf1", fmt.Sprintf("ssh,tcp,127.0.0.1,%d,,22", 2200 + i),
		"--nic2", "hostonly",
		"--hostonlyadapter2", adapter,
		"--cableconnected2", "on",
	    "--nicpromisc2", "allow-vms",
		).Run()
	if err == nil {
		log.Infof("Modified nic2 for \"%s\"", name)
	}
	return err
}

func Clone(baseName string, prefix string, num int, adapter string) error {
	for i := 1; i <= num; i++ {
		name := fmt.Sprintf("%s%03d", prefix, i)
		cmd := exec.Command(VBOX_MANAGE, "clonevm",
			baseName,
			"--snapshot", "origin",
			"--options", "link",
			"--name", name,
			"--register")
		out, err := cmd.Output()
		if err != nil {
			return err
		} else {
			err = Modify(name, adapter, i)
			log.Infof("Clone: %s", strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func Remove(args ...string) error {
	for _, name := range args {
		if name == "base" {
			err := exec.Command(VBOX_MANAGE, "snapshot", "base", "delete", "origin").Run()
			if err != nil {
				log.Info("Removed snapshot \"base/origin\"")
			}
		}
		cmd := exec.Command(VBOX_MANAGE, "unregistervm", name, "--delete")
		_, err := cmd.Output()
		log.Infof("Removed \"%s\"", name)
		if err != nil {
			return err
		}
	}
	return nil
}
