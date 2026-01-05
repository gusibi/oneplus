import os
import socket
import io
import qrcode
from fastapi import FastAPI, UploadFile, File, HTTPException
from fastapi.responses import HTMLResponse, StreamingResponse, FileResponse
from fastapi.staticfiles import StaticFiles
from fastapi.middleware.cors import CORSMiddleware
from typing import List
import uvicorn

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Setup directories
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
UPLOAD_DIR = os.path.join(BASE_DIR, "uploads")
STATIC_DIR = os.path.join(BASE_DIR, "static")

os.makedirs(UPLOAD_DIR, exist_ok=True)
os.makedirs(STATIC_DIR, exist_ok=True)

def get_local_ip():
    """Detect the local IP address of the machine."""
    try:
        # Connect to a public DNS server to determine the best local IP
        s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
        s.connect(("8.8.8.8", 80))
        ip = s.getsockname()[0]
        s.close()
        return ip
    except Exception:
        return "127.0.0.1"

LOCAL_IP = get_local_ip()
PORT = 8000
BASE_URL = f"http://{LOCAL_IP}:{PORT}"

# Mount static files
app.mount("/static", StaticFiles(directory=STATIC_DIR), name="static")

@app.get("/", response_class=HTMLResponse)
async def read_root():
    """Serve the single page application."""
    index_path = os.path.join(STATIC_DIR, "index.html")
    if os.path.exists(index_path):
        with open(index_path, "r", encoding="utf-8") as f:
            return f.read()
    return "<h1>Frontend not found. Please ensure static/index.html exists.</h1>"

@app.post("/api/upload")
async def upload_files(files: List[UploadFile] = File(...)):
    """Handle file uploads from the host machine."""
    saved_files = []
    for file in files:
        filename = file.filename
        # Simple sanitization to prevent directory traversal
        filename = os.path.basename(filename)
        file_location = os.path.join(UPLOAD_DIR, filename)
        
        with open(file_location, "wb+") as file_object:
            # Read in chunks to avoid memory issues with large files
            while content := await file.read(1024 * 1024):  # 1MB chunks
                file_object.write(content)
        saved_files.append(filename)
    return {"uploaded": saved_files}

@app.get("/api/files")
async def list_files():
    """List all available files for download."""
    files = []
    if os.path.exists(UPLOAD_DIR):
        files = [
            {"name": f, "size": os.path.getsize(os.path.join(UPLOAD_DIR, f))}
            for f in os.listdir(UPLOAD_DIR) 
            if os.path.isfile(os.path.join(UPLOAD_DIR, f)) and not f.startswith('.')
        ]
    return {"files": files}

@app.get("/api/download/{filename}")
async def download_file(filename: str):
    """Download a specific file."""
    # Security check: ensure filename is just a basename
    filename = os.path.basename(filename)
    file_path = os.path.join(UPLOAD_DIR, filename)
    
    if os.path.exists(file_path):
        return FileResponse(file_path, filename=filename)
    raise HTTPException(status_code=404, detail="File not found")

@app.get("/api/qr")
async def get_qr():
    """Generate QR code for the current local URL."""
    img = qrcode.make(BASE_URL)
    buf = io.BytesIO()
    img.save(buf, format="PNG")
    buf.seek(0)
    return StreamingResponse(buf, media_type="image/png")

@app.get("/api/info")
async def get_info():
    """Return server info."""
    return {"ip": LOCAL_IP, "port": PORT, "url": BASE_URL}

if __name__ == "__main__":
    print(f"Starting server at {BASE_URL}")
    uvicorn.run(app, host="0.0.0.0", port=PORT)
