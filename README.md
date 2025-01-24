# SSO-Auth-GO package

Simple package for authenticate user for pet-projects without having to start your own auth-service.

---

# Commands  
  

### Build

```shell
go build ./cmd/sso/main.go
```

### Run

```shell
go run ./cmd/sso/main.go --config="./path/to/config.yaml
```

---

# Flags

| flag   | example               | required | default             |
|--------|-----------------------|----------|---------------------|
| config | ./path/to/config.yaml | _false_  | ./config/local.yaml |