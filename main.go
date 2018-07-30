package main

import (
	"os"
	"fmt"
	"os/exec"
	"syscall"
	"path/filepath"
	"io/ioutil"
	"strconv"
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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
		Unshareflags: syscall.CLONE_NEWNS,
	}
	must(cmd.Run())
}

func child(){
	fmt.Printf("running %v as pid: %d\n", os.Args[2:], os.Getpid())
	cg()
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	must(cmd.Run())
}

func cg(){
	cgroups := "/sys/fs/cgroup/"
	pids := filepath.Join(cgroups,"pids")
	fmt.Println("Creating.. pids at: ",filepath.Join(pids,"chauhr"))

	must(os.MkdirAll(filepath.Join(pids,"chauhr"), 0755))
	must(ioutil.WriteFile(filepath.Join(pids, "chauhr/pids.max"), []byte("10"), 0700))
	must(ioutil.WriteFile(filepath.Join(pids,"chauhr/notify_on_release"), []byte("1"), 0700)) // release the cgroup when the container ends.
	must(ioutil.WriteFile(filepath.Join(pids,"chauhr/cgroup.procs"), []byte(strconv.Itoa(os.Getpid())),0700))
}

func must (err error){
	if err != nil{
		panic(err)
	}
}
