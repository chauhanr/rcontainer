package main

import (
	"os"
	"os/exec"
	"syscall"
)

func main(){
	switch os.Args[1]{
		case "run":
			run()
	  	default:
			panic("command not supported")
	}
}

func run(){
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=rcontainer #"}
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWNS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER ,
	}
	must(cmd.Run())
}

func must (err error){
	if err != nil{
		panic(err)
	}
}
