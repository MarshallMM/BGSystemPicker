@ECHO OFF
cd C:\Users\marsh\BGSystemPicker\src >NUL
C:\Users\marsh\BGSystemPicker\build.sh >NUL
ECHO Remember to stop service
echo.| plink.exe -t -i C:\Users\marsh\.ssh\puttykey.ppk pi@192.168.1.77 "sudo systemctl stop systemBot.service"
scp BGSystemPicker pi@192.168.1.77:/home/pi/systemBot
echo.| plink.exe -t -i C:\Users\marsh\.ssh\puttykey.ppk pi@192.168.1.77 "sudo systemctl start systemBot.service"
del BGSystemPicker