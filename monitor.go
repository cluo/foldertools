package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"
)

func walkFn(path string, info os.FileInfo, err error) error {
	fmt.Printf("%s \n", path)
	return nil
}

func main() {

	dir := flag.String("dir", "./", "目录地址")
	execSh := flag.String("e", "du -h {file}", "执行的命令,{file}为改变的文件")
	flag.Parse()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				//		log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					//		log.Println("modified file:", event.Name)
					execCommand(event.Name, *execSh)
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(*dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done

}

func execCommand(path string, execSh string) error {

	ext := path[len(path)-4:]

	if ext != ".swp" {
		execStr := strings.Replace(execSh, "{file}", path, -1)

		execSp := strings.Split(execStr, " ")
		cmd := exec.Command(execSp[0], execSp[1:]...)
		out, err := cmd.Output()

		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(out))
	}

	return nil
}
