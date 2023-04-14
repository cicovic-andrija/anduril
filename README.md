# Andúril notes server

## Deploy

```
ping notes.acicovic.me
ssh root@notes.acicovic.me
mkdir -p /srv/anduril
mkdir -p /srv/anduril/data
exit

rsync -v ./deploy/anduril.service root@notes.acicovic.me:/etc/systemd/system/
rsync -v ./deploy/anduril-config.json root@notes.acicovic.me:/srv/anduril/data/
rsync -rv static templates root@notes.acicovic.me:/srv/anduril/data/
rsync -v ./out/anduril-server root@notes.acicovic.me:/srv/anduril/

ssh root@notes.acicovic.me
systemctl enable anduril.service
systemctl start anduril.service
systemctl status anduril.service
exit
```
