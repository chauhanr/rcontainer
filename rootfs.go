package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
)

func pivotRoot(newroot string) error {
	putold := filepath.Join(newroot, "/.pivot_root")
	// the pivot root must get a new root and old root and both should be
	// in a different file system when compared to the newRoot.
	if err := syscall.Mount(newroot, newroot, "", syscall.MS_BIND|syscall.MS_REC, ""); err != nil {
		return err
	}
	// create a old directory
	if err := os.MkdirAll(putold, 0700); err != nil {
		return err
	}
	// call the pivot root function
	if err := syscall.PivotRoot(newroot, putold); err != nil {
		return err
	}

	if err := os.Chdir("/"); err != nil {
		return err
	}

	putold = "/.pivot_root"
	if err := syscall.Unmount(putold, syscall.MNT_DETACH); err != nil {
		return err
	}

	if err := os.RemoveAll(putold); err != nil {
		return err
	}
	return nil
}

func exitIfRootfsNotFound(rootfsPath string) {
	if _, err := os.Stat(rootfsPath); os.IsNotExist(err) {
		usefulMsg := fmt.Sprintf(`%s does not exist please create this directory use the busybox tar and unpack it on the system
	   and unpack it 
	   mkdir -p %s
	   tar -C %s -xf busybox.tar
	   `, rootfsPath, rootfsPath, rootfsPath)
		fmt.Println(usefulMsg)
		os.Exit(1)
	}

}

func mountProc(newroot string) error {
	source := "proc"
	target := filepath.Join(newroot, "/proc")
	fstype := "proc"
	flags := 0
	data := ""

	os.Mkdir(target, 0755)
	if err := syscall.Mount(
		source,
		target,
		fstype,
		uintptr(flags),
		data,
	); err != nil {
		return err
	}

	return nil
}
