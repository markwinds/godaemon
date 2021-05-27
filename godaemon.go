package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
)

func init() {
	goDaemon := flag.Bool("d", false, "run app as a daemon with -d=true.")
	kill := flag.Bool("k", false, "kill app which run in background")
	for _, arg := range os.Args[1:] {
		if arg == "-d" || arg == "-d=true" {
			*goDaemon = true
		}
		if arg == "-k" || arg == "-k=true" {
			*kill = true
		}
	}

	if *goDaemon {
		cmd := exec.Command(os.Args[0], flag.Args()...)
		if err := cmd.Start(); err != nil {
			fmt.Printf("start %s failed, error: %v\n", os.Args[0], err)
			os.Exit(1)
		}
		fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid)
		os.Exit(0)
	}

	if *kill {
		err := killByName(os.Args[0])
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}
}

func killByName(name string) error {
	r := regexp.MustCompile(`[/\\]`)
	s := r.Split(name, -1)
	name = s[len(s)-1]
	var cmd *exec.Cmd
	sys := runtime.GOOS
	switch sys {
	case "linux":
		cmd = exec.Command("pkill", name)
	case "windows":
		if !(len(name) > 4 && name[len(name)-4:] == ".exe") {
			name = name + ".exe"
		}
		cmd = exec.Command("taskkill", `/f`, `/im`, name, "/T") //T参数用来杀死其子进程
	default:
		return fmt.Errorf("-k arg just support windows and linux")
	}
	if err := cmd.Start(); err != nil {
		fmt.Printf("cmd [%s] start failed because:\n", cmd.String())
		return err
	}
	fmt.Printf("kill [%s] successful", name)
	return nil
}
