import socket
import threading
import time
import csv
import os
import uuid
from sessions import (
    lock, seshxtoken, seshxsocket,
    validsession, cleansession, appendcsv, csv_lock
)
from paths import USUARIOS_CSV, HISTORIAL_CSV
HOST = "0.0.0.0"
TCP_PORT = 9000
UDP_PORT = 9001


def csvexiste(path, header):
    if not os.path.exists(path):
        with open(path, 'w', newline='', encoding='utf-8') as f:
            csv.writer(f).writerow(header)


def credcheck(username, password):
    with csv_lock:
        with open(USUARIOS_CSV, 'r', newline='', encoding='utf-8') as f:
            reader = csv.reader(f)
            next(reader, None)
            for row in reader:
                if len(row) >= 2 and row[0] == username and row[1] == password:
                    return True
    return False


def appendhist(username, message):
    with csv_lock:
        with open(HISTORIAL_CSV, 'a', newline='', encoding='utf-8') as f:
            csv.writer(f).writerow([time.time(), username, message])


def broadcast(sender_username, message, exclude_conn):
    with lock:
        targets = [s["socket"] for s in seshxtoken.values() if s["socket"] != exclude_conn]
    line = f"INCOMING {sender_username} {message}\n".encode("utf-8")
    for sock in targets:
        try:
            sock.sendall(line)
        except OSError:
            pass


def handlelogin(conn, username, password):
    if not credcheck(username, password):
        conn.sendall(b"ERROR INVALID_CREDENTIALS\n")
        return
    token = uuid.uuid4().hex
    now = time.time()
    with lock:
        seshxtoken[token] = {
            "username": username,
            "socket": conn,
            "created": now,
            "last_heartbeat": None,
        }
        seshxsocket[conn] = token

    appendcsv(token, username, now)
    conn.sendall(f"OK {token} {UDP_PORT}\n".encode("utf-8"))

def handlemsg(conn, token, message):
    with lock:
        session = seshxtoken.get(token)
        valid = session is not None and validsession(session)
        username = session["username"] if session else None

    if not valid:
        if session is None:
            conn.sendall(b"ERROR SESSION_INVALID_TOKEN\n")
        else:
            conn.sendall(b"ERROR SESSION_EXPIRED\n")
        return

    appendhist(username, message)
    conn.sendall(b"ACK\n")
    broadcast(username, message, exclude_conn=conn)


def dispatch(conn, line):
    print(line)
    parts = line[0].strip().split(" ", 2)
    if not parts or parts[0] == "":
        return

    cmd = parts[0]

    if cmd == "LOGIN" and len(parts) == 3:
        params = line[-1].strip().split(" ")
        handlelogin(conn, params[0], params[1])
    elif cmd == "MSG" and len(parts) == 3:
        params = line[-1].strip().split(" ")
        handlemsg(conn, params[0], params[1])
    elif cmd == "LOGOUT":
        cleansession(conn=conn)
    else:
        conn.sendall(b"ERROR UNKNOWN_COMMAND\n")

def handle_client(conn, addr):
    buffer = ""
    try:
        while True:
            data = conn.recv(4096)
            if not data:
                break
            buffer += data.decode("utf-8", errors="replace")
            if "\n" in buffer:
                lines = buffer.strip().split("\n")

                request = []
                for i in range(len(lines)):
                    request.append(lines[i])

                if request:
                    dispatch(conn, request)
                buffer = ""
                
    except (ConnectionResetError, OSError):
        pass
    finally:
        cleansession(conn=conn)








def main():
    csvexiste(USUARIOS_CSV, ["username", "password", "fecha_registro"])
    csvexiste(HISTORIAL_CSV, ["timestamp", "username", "mensaje"])
    server_sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    server_sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
    server_sock.bind((HOST, TCP_PORT))
    server_sock.listen()
    print(f"TCP server listening on {HOST}:{TCP_PORT}")



    while True:
        conn, addr = server_sock.accept()
        t = threading.Thread(target=handle_client, args=(conn, addr), daemon=True)


        t.start()
if __name__ == "__main__":
    main()
