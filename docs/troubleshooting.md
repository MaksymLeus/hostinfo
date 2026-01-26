# Troubleshooting

## HELM
### Check Pod Logs
kubectl logs deployment/<name> -n <namespace>

### Check Pod Events
kubectl describe pod <pod> -n <namespace>

### Common Issues

| Issue | Cause | Fix |
|---|---|---|
| CrashLoopBackOff | Misconfigured env/secret | Check secret values |
| ImagePullBackOff | Bad image URL / auth | Check registry creds |
| Ingress 404 | Wrong host/path | Match ingress host/path |

