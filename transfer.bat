cd C:\Users\marsh\BGSystemPicker\src
#env GOOS=linux GOARCH=arm GOARM=5 go build
C:\Users\marsh\BGSystemPicker\build.sh
scp BGSystemPicker pi@192.168.1.77:/home/pi/systemBot
pause
del BGSystemPicker