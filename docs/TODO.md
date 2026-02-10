
# Upgrade UI
- Make new ui on js
- favicon.ico 


# CI
Scenario B — .releaserc stored in another repo

Option 1 — Checkout config repo first

In workflow:
```yaml
- uses: actions/checkout@v4
  with:
    path: current

- uses: actions/checkout@v4
  with:
    repository: myorg/semantic-config
    path: semantic-config

- name: Use external .releaserc
  run: cp semantic-config/.releaserc current/.releaserc
```

This approach is used in mono-repos and org-wide standards.