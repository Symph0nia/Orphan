# Orphan
## 0x01 简介

Orphan是一款非流量触发的Ring 3 用户态Linux后门，可以有效规避通过流量检测、端口检测的防护手段。

## 0x02 问题来源

常见的后门监听方式为，在某端口上开启长连接，例如反弹shell：

![image-20240501152314931](https://raw.githubusercontent.com/ZacharyZcR/Orphan/main/images/1.png)

会存在连接端口被发现的问题。

也有一些后门使用定时短连接的方式去取回待执行的命令，也存在定时的可疑流量问题。

有没有一种方式可以实现完全无流量的后门方式呢？

思路如下：

通过Nginx日志等文件方式，存储待执行的命令，后门程序监控Nginx日志，攻击者将待执行的命令混在流量中发送到Nginx，程序从Nginx日志获取待执行命令，从而实现RCE。

可惜Nginx默认日志为Root权限，用户权限无法读取。

根据360安全忍者的提示，在/proc目录下找寻可用功能点。

## 0x03 实现方式

Linux的/proc/net/sockstat 是用来监控系统流量的组件，内容如下：

```bash
sockets: used 209
TCP: inuse 20 orphan 0 tw 2 alloc 67 mem 5
UDP: inuse 3 mem 5
UDPLITE: inuse 0
RAW: inuse 0
FRAG: inuse 0 memory 0
```

其中orphan孤儿套接字是那些已经没有任何应用程序与之关联的套接字，通常是在等待关闭的状态。

通过trigger/orphan.py进行快速的流量发送

```bash
sockets: used 215
TCP: inuse 30 orphan 13 tw 2 alloc 76 mem 17
UDP: inuse 3 mem 5
UDPLITE: inuse 0
RAW: inuse 0
FRAG: inuse 0 memory 0
```

发现会导致orphan流量显著增加

根据以上原理进行程序设计，后门程序将每隔5S监控/proc/net/sockstat，如果orphan流量超过10，即开启一分钟的后门。

![image-20240501153505318](https://raw.githubusercontent.com/ZacharyZcR/Orphan/main/images/2.png)

![image-20240501153526910](https://raw.githubusercontent.com/ZacharyZcR/Orphan/main/images/3.png)

停止trigger/orphan.py后，后门自动关闭。

![image-20240501153615458](https://raw.githubusercontent.com/ZacharyZcR/Orphan/main/images/4.png)

等待下一次的启动。

## 0x04 进阶

理论上来说，更加详细的orphan流量设计，可以让/proc/net/sockstat的orphan显示呈现规律变化，可以传递0/1二进制编码，用于摩斯电码传递信息。

具体实现比较复杂，这里暂时不进行实现，仅提供思路。

通过有规律的TCP连接，可以将待执行的命令直接传输到后门上，实现0流量后门。

仅调整扫描器的频率即可实现RCE效果。
