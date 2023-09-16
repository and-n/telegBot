env GOOS=linux GOARCH=arm GOARM=5 go build
if [ $? -eq 0 ]; then
    echo "All done"
else
    echo "Error"
fi