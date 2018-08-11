package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/docker/pkg/reexec"
)

func main() {
	var rootfs string
	flag.StringVar(&rootfs, "rootfs", "/tmp/ns-process/rootfs", "Path to root file system to use")
	flag.Parse()

	exitIfRootfsNotFound(rootfs)

	cmd := reexec.Command("nsInitialization", rootfs)
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

func init() {
	reexec.Register("nsInitialization", nsInitialization)
	if reexec.Init() {
		os.Exit(0)
	}
}

func nsInitialization() {
	//fmt.Printf("\n>> namespace setup code goes here << \n\n")
	/*switch os.Args[1] {
	case "run":
		run()
	default:
		panic("command not supported")
	}*/
	newrootPath := os.Args[1]
	if err := mountProc(newrootPath); err != nil {
		fmt.Printf("Error mounting /proc - %s \n", err)
		os.Exit(1)
	}
	if err := pivotRoot(newrootPath); err != nil {
		fmt.Printf("Error using the pivot root- %s \n", err)
		os.Exit(1)
	}
	if err := syscall.Sethostname([]byte("rcontainer")); err != nil {
		fmt.Println("Error in setting the host name")
		os.Exit(1)
	}
	run()
}

func run() {

	cmd := exec.Command("/bin/sh")
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr

	cmd.Env = []string{"PS1=-[rcontainer]- # "}

	if err := cmd.Run(); err != nil {
		fmt.Printf("Error running the /bin/sh command - %s\n", err)
		os.Exit(1)
	}
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
