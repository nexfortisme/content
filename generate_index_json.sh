CONTENT_ROOT=./posts go run scripts/main.go

echo "Copying File to Blog Public Dir"
cp ./index.json ../blog/public/content