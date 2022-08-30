cd C:\Users\marsh\BGSystemPicker\bot
scp BGSystemPicker pi@192.168.1.75:/home/pi/systemBot
pause
//env GOOS=linux GOARCH=arm GOARM=5 go build