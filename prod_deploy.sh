#!/bin/bash
# download config file from firestore
wget -O config.yml https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/config.prod.yml\?alt\=media\&token\=00b9af03-be96-4a06-b80f-57be3014594c

# download swiftex-firebase.json from firestore
wget -O swiftex-firebase.json https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/swiftex-firebase.json\?alt\=media\&token\=d7ee69b0-0967-48ad-8bb9-9590c223aa94

cat wget-log

docker build -t caffeines/swiftex:1.0 .

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