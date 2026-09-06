import threading
import tcp_server
import udp_server
# Punto de entrada único que levanta ambos servidores (TCP y UDP) en el mismo proceso
if __name__ == "__main__":
    threading.Thread(target=udp_server.main, daemon=True).start()
    tcp_server.main()
