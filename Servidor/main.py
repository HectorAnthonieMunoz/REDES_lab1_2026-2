import threading
import tcp_server
import udp_server
if __name__ == "__main__":
    threading.Thread(target=udp_server.main, daemon=True).start()
    tcp_server.main()