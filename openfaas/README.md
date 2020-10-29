Sourcehawk OpenFaaS
-------------------

## Function URL
In the below examples `<function-url>` is the URL in which the function is exposed 
in an OpenFaas platform.  For example: `https://sourcehawk-scan.openfaas-example.com`

### Scanning a Public Repository
Below is an example of scanning a public repository in github

```bash
curl -X POST <function-url>/org/repo
```

### Scanning a Public Repository (With a specific Ref)
Below is an example of scanning a public repository in github with the abbreviated commit hash `894c000`

```bash
curl -X POST <function-url>/org/repo/894c000
```

Below is an example of scanning a public repository in github with the branch `develop`

```bash
curl -X POST <function-url>/org/repo/develop
```

### Scanning a Private Repository
Below is an example of scanning a private repository in github

```bash
GITHUB_PAT="abc123"
curl -X POST -H "Authorization: ${GITHUB_PAT}" <function-url>/org/repo
```

### Scanning a Repository in Github Enterprise
Below is an example of scanning a repository in github enterprise where the domain name is `github.example.com`

```bash
GITHUB_PAT="abc123"
curl -X POST -H "Github-API-URL: https://github.example.com/api/v3" <function-url>/org/repo
```

## Docker Builds

### Scan
```sh
./scan/docker-build.sh
```

### Validate Config
```sh
./validate-config/docker-build.sh
```