#!/bin/bash
# download config file from firestore
wget -O config.yml https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/config.prod.yml\?alt\=media\&token\=3e08bdea-abf9-479e-92c7-2bd4e2f9448c

# download swiftex-firebase.json from firestore
wget -O swiftex-firebase.json https://firebasestorage.googleapis.com/v0/b/swiftexbd.appspot.com/o/swiftex-firebase.json\?alt\=media\&token\=d7ee69b0-0967-48ad-8bb9-9590c223aa94

cat wget-log

docker build -t caffeines/swiftex:1.0.5 .

docker system prune