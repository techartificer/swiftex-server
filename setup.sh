#!/bin/bash
# download config file from firestore
wget -O config.yml https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/config.prod.yml\?alt\=media\&token\=b4b0a4ad-5b9f-41bd-8ca2-be2efa57950c

# download swiftex-firebase.json from firestore
wget -O swiftex-firebase.json https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/swiftex-firebase.json\?alt\=media\&token\=91db4d89-de7e-441c-8c79-2267656c0d88

cat wget-log

docker build -t swiftex:1.0 .

docker-compose down
docker-compose up -d

sleep 15
echo "Health check"
echo ""
curl -s localhost:4141 | json_pp
echo ""
echo ""
echo "Server is ready..."

docker system prune