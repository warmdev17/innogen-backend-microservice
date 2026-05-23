# VPS Deployment with Nginx + HTTPS

## Architecture

```
Internet → Nginx (HTTPS :443) → 127.0.0.1:8080 (api_gateway)
                                   → Docker network → backend services
```

Only api_gateway is exposed. All other services stay private inside Docker.

## Prerequisites

- Ubuntu 22.04+ VPS
- Docker + Docker Compose
- Domain pointed to VPS IP (e.g., api.maiphuongtrunghieu.site)

## 1. Clone and set up

```bash
git clone git@github.com:warmdev17/innogen-backend-microservice.git
cd innogen-backend-microservice

# Create secrets directory
mkdir -p secrets
# Place GitHub App private key:
# cp /path/to/private-key.pem secrets/github-key.pem
chmod 600 secrets/github-key.pem
```

## 2. Environment file

```bash
cp .env.example .env
# Edit .env with production values:
#   JWT_SECRET=<random-64-char>
#   GITHUB_APP_ID=3721867
#   GITHUB_PRIVATE_KEY_PATH=/app/secrets/github-key.pem
#   GITHUB_WEBHOOK_SECRET=<from-github-app-settings>
#   GITHUB_OAUTH_CLIENT_ID=Ov23li...
#   GITHUB_OAUTH_CLIENT_SECRET=<secret>
#   GITHUB_OAUTH_REDIRECT_URL=https://api.your-domain.com/github/oauth/callback
#   GITHUB_OAUTH_FRONTEND_REDIRECT_URL=https://your-domain.com/github
```

## 3. Start services

```bash
# Production (only gateway exposed)
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Apply schema + seed
docker compose exec -T postgres psql -U innogen -d innogen < schema.sql
docker compose exec -T postgres psql -U innogen -d innogen < seeds/dev_seed.sql
docker compose exec -T postgres psql -U innogen -d innogen < seeds/dev_problem_pack.sql
```

## 4. Install Nginx

```bash
sudo apt update && sudo apt install nginx -y
```

## 5. Nginx config

Create `/etc/nginx/sites-available/rinnogen`:

```nginx
server {
    listen 80;
    server_name api.maiphuongtrunghieu.site;

    client_max_body_size 2M;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 60s;
    }

    location /webhooks/github {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_read_timeout 30s;
    }
}
```

Enable:

```bash
sudo ln -sf /etc/nginx/sites-available/rinnogen /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 6. HTTPS with Certbot

```bash
sudo apt install certbot python3-certbot-nginx -y
sudo certbot --nginx -d api.maiphuongtrunghieu.site
```

Auto-renewal:

```bash
sudo certbot renew --dry-run  # test
# Auto-renews via systemd timer
```

## 7. GitHub App Production URLs

After HTTPS is set up, update GitHub App settings:
- **Callback URL**: `https://api.your-domain.com/github/oauth/callback`
- **Webhook URL**: `https://api.your-domain.com/webhooks/github`
- **Setup URL**: `https://your-domain.com/github`

## 8. Firewall

```bash
sudo ufw allow 22    # SSH
sudo ufw allow 80    # HTTP (for certbot)
sudo ufw allow 443   # HTTPS
sudo ufw deny 8080   # Block direct gateway access
sudo ufw deny 8081   # Block auth_service
sudo ufw deny 8082   # Block run_service
sudo ufw deny 8083   # Block submission_service
sudo ufw deny 8084   # Block repo_service
sudo ufw deny 5432   # Block postgres
sudo ufw deny 6379   # Block redis
sudo ufw deny 2000   # Block piston
sudo ufw enable
```

## 9. Secrets mount

In `docker-compose.prod.yml`, add to repo_service:

```yaml
  repo_service:
    volumes:
      - ./secrets:/app/secrets:ro
```

Set `GITHUB_PRIVATE_KEY_PATH=/app/secrets/github-key.pem`

## 10. Useful commands

```bash
# View logs
docker compose logs -f api_gateway

# Restart after code change
docker compose -f docker-compose.yml -f docker-compose.prod.yml up -d --build

# Database backup
docker compose exec postgres pg_dump -U innogen innogen > backup.sql

# Check health
curl -s https://api.your-domain.com/health
```

## Production Checklist

- [ ] .env secrets set (JWT, GitHub, OAuth)
- [ ] GitHub private key in secrets/
- [ ] Nginx HTTPS working
- [ ] Firewall active
- [ ] GitHub App callback/webhook URLs updated
- [ ] Health check passes
- [ ] Submission flow tested
- [ ] GitHub commit flow tested
