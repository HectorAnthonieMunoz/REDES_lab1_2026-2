import socket
import time
TOKEN = "9d40cfb700974c8d96ecc8673625e3ec"
s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
while True:
    s.sendto(f"HEARTBEAT {TOKEN}\n".encode(), ("localhost", 9001))
    print("sent heartbeat")
    time.sleep(3)