import threading
import time
import csv
from paths import SESIONES_CSV



lock = threading.Lock()
seshxtoken = {}
seshxsocket = {}
csv_lock = threading.Lock()
TOKEN_TTL = 600
FIRST_HEARTBEAT_GRACE = 30
HEARTBEAT_TIMEOUT = 60


def appendcsv(token, username, created):
    with csv_lock:
        with open(SESIONES_CSV, 'a', newline='', encoding='utf-8') as f:
            csv.writer(f).writerow([token, username, created, created, "ACTIVO"])
def writecsv():
    with csv_lock:
        with open(SESIONES_CSV, 'w', newline='', encoding='utf-8') as f:
            writer = csv.writer(f)
            writer.writerow(["token", "username", "timestamp_creacion", "timestamp_ultimo_heartbeat", "estado"])
            with lock:
                for token, s in seshxtoken.items():
                    writer.writerow([token, s["username"], s["created"], s["last_heartbeat"] or s["created"], "ACTIVO"])

def validsession(session):
    now = time.time()
    if now - session["created"] > TOKEN_TTL:
        return False
    if session["last_heartbeat"] is None:
        return (now - session["created"]) <= FIRST_HEARTBEAT_GRACE
    if now - session["last_heartbeat"] > HEARTBEAT_TIMEOUT:
        return False
    return True



def cleansession(conn=None, token=None):
    with lock:
        if conn is not None and token is None:
            token = seshxsocket.get(conn)
        if token is None:
            return
        seshxtoken.pop(token, None)
        sock = seshxsocket.pop(conn, None) if conn else None
        if sock is None:
            for c, t in list(seshxsocket.items()):
                if t == token:
                    sock = seshxsocket.pop(c, None)
                    conn = c
                    break
    if conn is not None:
        try:
            conn.close()
        except OSError:
            pass
    writecsv()