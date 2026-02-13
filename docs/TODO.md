
# Upgrade Frontend
- inject a random Party Parrot into your frontend every 1–5 minutes, like a fun Easter egg.

# Upgrade Backend
-

# CI

 — `.releaserc` stored in another repo

    Option — Checkout config repo first

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