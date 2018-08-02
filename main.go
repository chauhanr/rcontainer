package main

import (
	"os"
	"os/exec"
	"syscall"
	"fmt"
	"github.com/docker/docker/pkg/reexec"
)

func main(){
	cmd := reexec.Command("nsInitialization")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getuid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getgid(),
				Size:        1,
			},
		},
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the reexec.Command - %s\n", err)
		os.Exit(1)
	}
}

func init (){
	reexec.Register("nsInitialization", nsInitialization)
	if reexec.Init(){
		os.Exit(0)
	}
}

func nsInitialization(){
	fmt.Printf("\n>> namespace setup code goes here << \n\n")
	/*switch os.Args[1] {
	case "run":
		run()
	default:
		panic("command not supported")
	}*/
	run()
}

func run(){
	//cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd := exec.Command("/bin/sh")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[ns-process]- # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func must (err error){
	if err != nil{
		panic(err)
	}
}
