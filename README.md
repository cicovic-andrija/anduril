# How To Build and Deploy This Web Server

> To be used by the author and administrator of the website.

## Required Tools

- `go` version `1.20` or greater
- `git`
- `make`
- `openssl`
- `pandoc`
- `rsync`

**Important note: The server is designed to be run on a Linux operating system.**

## Build Process

The `Makefile` file in the root directory configures the `make` build system for all common operations. Execute `make` or `make build` to build the web server. Execute `make tools` to build the accompanying tools, and `make all` to build everything.

## Local (Dev) Deployment Process

1. Review and fill empty fields (passwords) where appropriate in the config file located in `configuration/`.
2. Execute `make devenv` to build everything, create a local working directory `out/` and copy necessary files to the working directory.
3. From the context of the `out/` directory, execute `./anduril-server &` to run the server as a background process.
4. Note: The server is configured with the `dev` configuration profile which allows only `localhost` connections on port `8080`, and presents a self-signed HTTPS certificate.

## Production Deployment Process

1. Review and fill empty fields (passwords) where appropriate in the config file located in `configuration/`.
2. Execute `make all` to build the server, accompanying tools and configuration files.
3. In case of first-time deployments, the following requirements must be met on the remote server machine (`notes.acicovic.me` host):
    1. Make sure `systemd`, `openssl`, `pandoc`, and `rsync` are installed on the system.
    2. Install the HTTPS certificate to the location indicated by the `https.network` section of the server's config file (e.g. instructions [https://letsencrypt.org/](https://letsencrypt.org/)).
    3. Add the needed SSH public key fingerprints to the SSH `known_hosts` file.
    4. From the local machine, send the `systemd` service config file to the remote machine: `rsync -v ./configuration/anduril.service {username}@notes.acicovic.me:/etc/systemd/system/`.
4. Sync all required working files to the remote machine with `rsync`, as described in "How To Sync Files to the Remote Machine".
5. Connect to the remote machine: `ssh {username}@notes.acicovic.me`.
6. In case of first time deployments, execute `systemctl enable anduril.service` followed by `systemctl start anduril.service`; otherwise simply run `systemctl restart anduril.service`.
7. Check the health status of the service: `systemctl status anduril.service`.
8. Monitor the main log file for errors for the first 10 minutes of execution.
9. Visit [https://notes.acicovic.me](https://notes.acicovic.me) in a web browser to confirm a new version has been successfully deployed.

# Working Directory Structure

| Location in repository | Location in working directory `/srv/anduril` |
| ---------------------- | ----------------------------- |
| `out/anduril-server` | `anduril-server` |
| `out/data/anduril-config.json` | `data/anduril-config.json` |
| `out/data/encrypted-config.txt` | `data/encrypted-config.txt` |
| `assets/templates/*.html` | `data/templates/*.html` |
| `assets/scripts/*.js` | `data/assets/*.js` |
| `assets/stylesheets/*.css` | `data/assets/*.css` |
| `assets/icons/*` | `data/assets/icons/*` |

# How To Sync Files to the Remote Machine

```
ping notes.acicovic.me
ssh root@notes.acicovic.me
mkdir -p /srv/anduril/data/assets
exit
rsync -v ./out/anduril-server root@notes.acicovic.me:/srv/anduril/
rsync -v ./out/data/encrypted-config.txt root@notes.acicovic.me:/srv/anduril/data/
rsync -rv assets/templates root@notes.acicovic.me:/srv/anduril/data/
rsync -v assets/scripts/*.js root@notes.acicovic.me:/srv/anduril/data/assets/
rsync -v assets/stylesheets/*.css root@notes.acicovic.me:/srv/anduril/data/assets/
rsync -rv assets/icons root@notes.acicovic.me:/srv/anduril/data/assets/
```

# Local Log Monitoring

```
ping notes.acicovic.me
ssh root@notes.acicovic.me
tail -f logs/* | cut -d$'\t' -f4 # OR $ tail -f logs/*
```
