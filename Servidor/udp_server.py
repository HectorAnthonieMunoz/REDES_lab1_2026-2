import socket
import threading
import time

from sessions import lock, seshxtoken, validsession, cleansession, writecsv

HOST = "0.0.0.0" 
UDP_PORT = 9001   # Puerto UDP donde se reciben los datagramas HEARTBEAT
WATCHDOG_INTERVAL = 5   # Cada cuántos segundos el watchdog revisa el estado de las sesiones


# Recibe datagramas HEARTBEAT <token> y actualiza el último heartbeat de la sesión correspondiente
def handlehb(sock):   
    while True:
        try:
            data, addr = sock.recvfrom(1024)
        except OSError:
            break
        msg = data.decode("utf-8", errors="replace").strip()
        parts = msg.split(" ", 1)
        if len(parts) != 2 or parts[0] != "HEARTBEAT":
            continue
        token = parts[1]
        found = False
        with lock:
            session = seshxtoken.get(token)
            if session:
                session["last_heartbeat"] = time.time()
                found = True
        if found:
            writecsv()



# Rutina en segundo plano que revisa periódicamente las sesiones y revoca las que expiraron
def watchdog():
    while True:
        time.sleep(WATCHDOG_INTERVAL)
        with lock:
            expired_tokens = [t for t, s in seshxtoken.items() if not validsession(s)]
        for token in expired_tokens:
            print(f"Watchdog: revoking expired session {token}")
            cleansession(token=token)

# Levanta el socket UDP, inicia el watchdog en un hilo aparte y comienza a recibir heartbeats
def main():
    sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    sock.bind((HOST, UDP_PORT))
    print(f"UDP server listening on {HOST}:{UDP_PORT}")

    threading.Thread(target=watchdog, daemon=True).start()
    handlehb(sock)


if __name__ == "__main__":
    main()
