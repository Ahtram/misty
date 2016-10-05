This is the private robot for projectA discord channel written by Golang.

This is just a note text file for me to record things that's not so important so please ignore it.

Server 連線設定
GoCrazy Server: 128.199.246.158 以 DDrd4y2EVbZY.ppk ssh 遠端登入 root
Create user: https://www.digitalocean.com/community/tutorials/how-to-create-a-sudo-user-on-ubuntu-quickstart
Add ssh key for user: https://www.digitalocean.com/community/questions/ssh-new-user-ubuntu-14-04
建立 user neo 並且產生 DDnd4y2EVbZY.ppk

在Server上安裝Go環境: https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04

建立 ssh key 讓 Server 可以 git clone 我們的 repo. 當然還須要改一下設定檔 config 讓 gitlab 直接使用該 private key

misty 專案 depend 下面的東西:
DiscordGo: https://github.com/bwmarrin/discordgo
[go get github.com/bwmarrin/discordgo] 
[go install -a github.com/bwmarrin/discordgo]
Color: https://github.com/fatih/color
[go get github.com/fatih/color] 
[go install -a github.com/fatih/color]
