package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/Sirupsen/logrus"
)

var VBOX_MANAGE = "VBoxManage"

func Import(file string, memory int, adapter string) error {
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

	err = Modify(name, adapter)

	_, err = exec.Command(VBOX_MANAGE, "snapshot", "base", "take", "origin").Output()
	log.Info("Snapshot \"origin\" taken")
	return err
}

func Modify(name string, adapter string) error {
	err := exec.Command(VBOX_MANAGE, "modifyvm", name,
		"--nic2", "hostonly",
		"--hostonlyadapter2", adapter,
		).Run()
	if err == nil {
		log.Infof("Modified nic2 for \"%s\"", name)
	}
	return err
}

func Clone(baseName string, prefix string, num int) error {
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
