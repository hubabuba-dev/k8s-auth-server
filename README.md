## TODO List


- [ ] Add configuration
- [ ] Add proper logging
- [ ] Add helmcharts
- [ ] Add tests

## Deployment
```dockerfile
# Build Docker image from Dockerfile and push to local registry
docker build -t localhost:5000/your-app:latest .
docker push localhost:5000/your-app:latest
```

## Manifests

### Auth
- `auth.yaml` - apply it in kube-apiserver manifest

### Certificates
- `certs.yaml` - creates certs in kubernetes
- `cert_exp_job.yaml` - download certs to local machine

### Roles
- contains basic admin role and rolebinding for it

### OIDC
- `depl.yaml` - app's deployment
- `ingress` - creates ingress for app
- `serv.yaml` - creates service for app
