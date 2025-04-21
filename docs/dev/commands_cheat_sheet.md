# **AI GENERATED**

# 🐳 Docker + Compose Cheatsheet for Infinite Experiment

## 🚀 Application Build & Run

### 🛠 Development (with `air`)
```bash
# Build and run dev stack
docker-compose -f docker-compose.dev.yml up --build

# Stop
docker-compose -f docker-compose.dev.yml down
```

### 🧱 Production
```bash
# Build and run prod stack
docker-compose -f docker-compose.prod.yml up -d --build

# Stop
docker-compose -f docker-compose.prod.yml down

# Rebuild only the API service
docker-compose -f docker-compose.prod.yml up -d --build api
```

---

## 🗃️ PostgreSQL Access

### 🔧 From Inside Container
```bash
# Open psql shell inside db container
docker exec -it $(docker ps -qf "name=db") psql -U ieuser -d infinite
```

### 💻 From Host (with Postgres installed)
```bash
psql -h localhost -p 5432 -U ieuser -d infinite
# Password: iepass
```

### 🌐 pgAdmin (Dev Only)
- Visit: http://localhost:5050
- Login: `admin@admin.com` / `admin`
- Add server:
  - Host: `db`
  - Port: `5432`
  - Username: `ieuser`
  - DB name: `infinite`

---

## 🧰 Handy Docker Commands

### 📦 Containers
```bash
# List running containers
docker ps

# Stop a container
docker stop <container_id>

# Remove a container
docker rm <container_id>
```

### 🐳 Images
```bash
# List images
docker images

# Remove an image
docker rmi <image_id>
```

### 🛠 Volumes
```bash
# List volumes
docker volume ls

# Inspect volume
docker volume inspect pgdata-dev
```

### 🧼 Cleanup
```bash
# Stop and remove all containers
docker stop $(docker ps -aq)
docker rm $(docker ps -aq)

# Remove all unused data
docker system prune -a
```

---

✅ Tip: Use `--build` whenever you change code or config and need fresh builds.

Let me know if you want to include reverse proxy, TLS cert, or monitoring commands too!
