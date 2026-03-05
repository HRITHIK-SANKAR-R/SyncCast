# SyncCast

# --- BINARIES & EXECUTABLES ---
# This keeps the compiled Go programs out of GitHub
synccast-lite
*.exe
*.exe~
*.dll
*.so
*.dylib

# --- MEDIA FILES (CRITICAL) ---
# Your team will be testing with movies. DO NOT upload them.
*.mp4
*.mkv
*.avi
*.mov
*.mp3
*.flac

# --- GO SPECIFIC ---
# Ignore dependencies (they are managed by go.mod)
vendor/
go.sum

# --- IDE & SYSTEM FILES ---
# Keeps the repo clean of your personal settings
.vscode/
.idea/
.DS_Store
Thumbs.db

# --- ENVIRONMENT & SECRETS ---
# If you eventually add API keys or secret tokens
.env
