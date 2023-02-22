# EDU-CLOUD API

Automation list
- [x] Build Dockerfile
- [x] Build Docker-compose
- [ ] Make github workflow

## Step to run
### Docker
Step 1 : build image
```
docker build -t ce-cloud-api .
```

Step: 2 run
```
docker run --rm -d -p 3001:3001 ce-cloud-api:latest
```

### Docker-compose
Step 1 : build image
```
docker-compose build 
```

Step: 2 run
```
docker-compose up -d
```