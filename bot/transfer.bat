cd C:\Users\marsh\BGSystemPicker\bot
#env GOOS=linux GOARCH=arm GOARM=5 go build
build.sh
scp BGSystemPicker pi@192.168.1.77:/home/pi/systemBot
pause
del BGSystemPicker