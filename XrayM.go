package main

import (
	"os/exec"
)

var XrayCmd *exec.Cmd

func Start() {
	XrayCmd = exec.Command("data/Xray/Xray", "start")
	XrayCmd.Start()
}

func Stop() {

}

func GetStatus() {

}
