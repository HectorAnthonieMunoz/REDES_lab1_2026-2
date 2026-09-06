import os
#Modula centralizado de rutas para los csv
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
USUARIOS_CSV = os.path.join(BASE_DIR, "..", "usuarios.csv")
SESIONES_CSV = os.path.join(BASE_DIR, "..", "sesiones.csv")
HISTORIAL_CSV = os.path.join(BASE_DIR, "..", "historial.csv")
