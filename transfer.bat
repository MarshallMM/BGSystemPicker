@ECHO OFF
cd C:\Users\marsh\BGSystemPicker\src >NUL
C:\Users\marsh\BGSystemPicker\build.sh >NUL
ECHO Remember to stop service
echo.| plink.exe -t -i C:\Users\marsh\.ssh\puttykey.ppk marsh@192.168.1.77 "sudo systemctl stop systemBot.service"
scp BGSystemPicker marsh@192.168.1.77:/home/marsh/systemBot
echo.| plink.exe -t -i C:\Users\marsh\.ssh\puttykey.ppk marsh@192.168.1.77 "sudo systemctl start systemBot.service"
