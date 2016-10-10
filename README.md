This is the private robot for projectA discord channel written by Golang.

This is just a note text file for me to record things that's not so important so please ignore it.

Server 連線設定教學
GoCrazy Server: 128.199.246.158 以 DDrd4y2EVbZY.ppk ssh 遠端登入 root
Create user: https://www.digitalocean.com/community/tutorials/how-to-create-a-sudo-user-on-ubuntu-quickstart
Add ssh key for user: https://www.digitalocean.com/community/questions/ssh-new-user-ubuntu-14-04
建立 user neo 並且產生 DDnd4y2EVbZY.ppk

在Server上安裝Go環境: https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04

建立 ssh key [GVSj1XnUGIT7c59] 讓 Server 可以 git clone 我們的 repo. 當然還須要改一下設定檔 config 讓 gitlab 直接使用該 private key

misty 專案 depend 下面的東西:
DiscordGo: https://github.com/bwmarrin/discordgo
[go get github.com/bwmarrin/discordgo] 
Color: https://github.com/fatih/color
[go get github.com/fatih/color] 

//讓 misty 加入自己 server 用的連結
https://discordapp.com/oauth2/authorize?client_id=231105148074721280&scope=bot&permissions=0

./misty -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.CtZU3A.vxmIdJICDZizVjAmM8908VrAD5c"
or
./misty -e [email] -p [password] -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.CtZU3A.vxmIdJICDZizVjAmM8908VrAD5c"
or 
./misty -e [email] -p [password]

使用 tmux 指令來背景常駐運作 misty

教學: https://gist.github.com/MohamedAlaa/2961058