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

	// 命令行k参数 用于杀死进程
	kill := flag.Bool("k", false, "kill app which run in background")

	// 命令行w参数 用于到监测程序退出值为10时重启程序 默认开启 使用-w禁用
	watcher := flag.Bool("w", true, "watch process and restart program when process exit 10")

	for _, arg := range os.Args[1:] {
		if arg == "-k" || arg == "-k=true" {
			*kill = true
		}
		if arg == "-w" || arg == "-w=true" {
			*watcher = false
		}
	}

	if *kill {
		err := killByName(os.Args[0])
		if err != nil {
			panic(err)
		}
		os.Exit(0)
	}

	if *watcher {
		// 重新启动程序 新进程需要-w参数 从而直接运行程序而不是启动watcher
		outProcess := false
		for !outProcess {
			newArgs := append(os.Args[1:], "-w")
			cmd := exec.Command(os.Args[0], newArgs...)
			cmd.Stdout = os.Stdout // 重定向输出到父进程的标准输出
			if err := cmd.Run(); err != nil {
				if exitError, ok := err.(*exec.ExitError); ok {
					if exitError.ExitCode() == 10 {
						continue // 当程序退出值为10时重启启动程序
					}
				}
			}
			outProcess = true
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
