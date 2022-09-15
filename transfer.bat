@ECHO OFF
cd C:\Users\marsh\BGSystemPicker\src >NUL
C:\Users\marsh\BGSystemPicker\build.sh >NUL
ECHO Remember to stop service
scp BGSystemPicker pi@192.168.1.77:/home/pi/systemBot
pause
del BGSystemPicker