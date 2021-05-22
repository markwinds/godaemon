package godaemon

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
	flag.Parse()

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
		cmd = exec.Command("taskkill", `/f`, `/im`, name)
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
