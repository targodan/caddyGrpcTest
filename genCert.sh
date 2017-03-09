#!/bin/bash

openssl req -x509 -newkey rsa:2048 -nodes -out server.crt -keyout server.key -days 365 -config <(cat <<-EOF
[ req ]
default_bits = 2048
prompt = no
default_md = sha256
x509_extensions = req_ext
distinguished_name = dn

[ dn ]
C=DE
ST=Some-State
L=More-Inof
O=Stuff
OU=Test
emailAddress=some@mail.de
CN=localhost

[ req_ext ]
subjectAltName = @alt_names

[ alt_names]
DNS.1 = localhost
DNS.2 = localhost.localdomain
IP = 127.0.0.1
EOF
)

