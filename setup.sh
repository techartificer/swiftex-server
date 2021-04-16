#!/bin/bash

rm wget-log

# download config file from firestore
wget -O config.yml -b https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/config.yml\?alt\=media\&token\=32652e75-6243-4666-8f1e-e15ac8d049a3

# download swiftex-firebase.json from firestore
wget -O swiftex-firebase.json -b https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/swiftex-firebase.json\?alt\=media\&token\=91db4d89-de7e-441c-8c79-2267656c0d88

cat wget-log

docker build -t swiftex .

docker-compose down
docker-compose up -d

sleep 15
echo "Health check"
echo ""
curl -s localhost:4141 | json_pp
echo ""
echo ""
echo "Server is ready..."