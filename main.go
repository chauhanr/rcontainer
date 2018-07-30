package main

import (
	"os"
	"fmt"
	"os/exec"
	"syscall"
)

func main(){
	switch os.Args[1]{
		case "run":
			run()
	    case "child":
			child()
		default:
			panic("command not supported")
	}
}

func run(){
	cmd := exec.Command("/proc/self/exe", append([]string{"child"},os.Args[2:]...)...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
	must(cmd.Run())
}

func child(){
	fmt.Printf("running %v as pid: %d\n", os.Args[2:], os.Getpid())
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func must (err error){
	if err != nil{
		panic(err)
	}
}
