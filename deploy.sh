scp ./telegBot pi@192.168.50.60:go_bot/updates/telegBot
if [ $? -eq 0 ]; then
    echo "File successfully copied to the server."
else
    echo "Error: File copy failed."
fi
