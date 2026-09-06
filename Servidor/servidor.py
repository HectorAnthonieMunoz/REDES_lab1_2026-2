from http.server import HTTPServer, BaseHTTPRequestHandler
import json
import csv
import os
import threading
from datetime import datetime
from paths import USUARIOS_CSV, HISTORIAL_CSV
PORT = 8888 
HOST = "localhost"
lock = threading.Lock() # Lock para proteger el acceso concurrente a usuarios.csv
# Crea el archivo CSV con su encabezado si todavía no existe
def csvexiste(path, header):
    if not os.path.exists(path):
        with open(path, 'w', newline='', encoding='utf-8') as f:
            csv.writer(f).writerow(header)


# Revisa si un username ya está registrado en usuarios.csv
def usuarioexiste(username):
    with open(USUARIOS_CSV, 'r', newline='', encoding='utf-8') as f:
        reader = csv.reader(f)
        next(reader, None)
        for row in reader:
            if row and row[0] == username:
                return True
    return False



# Handler HTTP que implementa los endpoints POST /register y GET /history
class servidor(BaseHTTPRequestHandler):
    def do_GET(self):
        if self.path == "/history":
            rows = []
            with open(HISTORIAL_CSV, 'r', newline='', encoding='utf-8') as f:
                reader = csv.reader(f)
                next(reader, None)
                for row in reader:
                    rows.append(",".join(row))
            body = "\n".join(rows).encode("utf-8")
            self.send_response(200)
            self.send_header("Content-Type", "text/plain; charset=utf-8")
            self.send_header("Content-Length", str(len(body)))
            self.end_headers()
            self.wfile.write(body)
        else:
            self.send_error(404)

    def do_POST(self):
        if self.path == "/register":
            length = int(self.headers.get("Content-Length", 0))
            line = self.rfile.read(length).decode("utf-8")
            tipo = self.headers.get("Content-Type", "")

            if "json" not in tipo:
                self.send_error(400)
                return
            try:
                datos = json.loads(line)
            except json.JSONDecodeError:
                self.send_error(400)
                return


            user = datos.get("username")
            passwrd = datos.get("password")

            if not user or not passwrd:
                self.send_error(400)
                return
            
            with lock:            
                if usuarioexiste(user):
                    self.send_error(409)
                    return

                fecha_registro = datetime.now().isoformat()
                with open(USUARIOS_CSV, 'a', newline='', encoding='utf-8') as f:
                    csv.writer(f).writerow([user, passwrd, fecha_registro])

            self.send_response(201)
            self.end_headers()
        else:
            self.send_error(404)




# Inicializa los CSV necesarios y arranca el servidor HTTP  
def main():
    csvexiste(USUARIOS_CSV, ["username", "password", "fecha_registro"])
    csvexiste(HISTORIAL_CSV, ["timestamp", "username", "mensaje"])
    runserver = HTTPServer((HOST, PORT), servidor)
    runserver.serve_forever()
if __name__ == "__main__":
    main()
