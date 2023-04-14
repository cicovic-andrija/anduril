#!/bin/bash
openssl req \
    -x509 \
    -newkey rsa:4096 \
    -sha256 \
    -days 365 \
    -nodes \
    -out ./out/certfile.crt \
    -keyout ./out/keyfile.key \
    -subj "/C=RS/ST=Belgrade/L=Belgrade/O=Notes Website/OU=Personal/CN=localhost"
