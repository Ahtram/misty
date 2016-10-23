This is the private robot for projectA discord channel written by Golang.

[目標]
1. Misty 可以讓玩家查詢遊戲內的煉金合成公式。
2. 公式資料來源為 Google Sheet 內的資料，所以 Bot 回傳的內容會隨著文件更新一併改變，目前預定作法為非即時查詢，這樣比較省效能。
3. 支援多語系(繁簡英日)。
4. 其他功能想到的話可以後續再繼續擴充。 

[Todo]
[ok] 讓 bot 可以接受 misty 這個關鍵字作為命令開頭。（只對 misty 開頭的字串做判斷處理，這邊得先做出字串用空白分割)
[ok] 讀取設定sheet，來決定接收命令時要回覆什麼特定字串。(一樣用 Google Sheet 定義)
[ok] 新增 update command 讓 misty 可以由任意使用者下指令更新資料。
支援 Beam 或者 hitbox (特定使用者頻道) 開始直播的 Discord 通知。
測試 golong timer 以及隨機自言自語功能的可能性。

--

Below is just some random note text for me to remember things that's not so important so please ignore it.

Server 連線設定教學
GoCrazy Server: 128.199.246.158 以 DDrd4y2EVbZY.ppk ssh 遠端登入 root
Create user: https://www.digitalocean.com/community/tutorials/how-to-create-a-sudo-user-on-ubuntu-quickstart
Add ssh key for user: https://www.digitalocean.com/community/questions/ssh-new-user-ubuntu-14-04
建立 user neo 並且產生 DDnd4y2EVbZY.ppk

在Server上安裝Go環境: https://www.digitalocean.com/community/tutorials/how-to-install-go-1-6-on-ubuntu-14-04

建立 ssh key [GVSj1XnUGIT7c59] 讓 Server 可以 git clone 我們的 repo. 當然還須要改一下設定檔 config 讓 gitlab 直接使用該 private key

misty 專案 depend 下面的 package :
DiscordGo: https://github.com/bwmarrin/discordgo
[go get github.com/bwmarrin/discordgo] 
Color: https://github.com/fatih/color
[go get github.com/fatih/color] 

Discord Application 後台
https://discordapp.com/developers/applications/me/231105148074721280

//讓 misty 加入自己 server 用的連結
https://discordapp.com/oauth2/authorize?client_id=231105148074721280&scope=bot&permissions=0

//執行指令
./misty -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.Cus7Sw.n-suc_aXypKw-EnkRw8kA3TMU4Y"
or
./misty -e [email] -p [password] -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.Cus7Sw.n-suc_aXypKw-EnkRw8kA3TMU4Y"
or 
./misty -e [email] -p [password]

使用 tmux 指令來背景常駐運作 misty
tmux 教學: https://gist.github.com/MohamedAlaa/2961058

擴充命令用 Google Sheet
https://docs.google.com/spreadsheets/d/1haLbQuE7TtF79_J2XLbzFRYbAkfGRCmrXxwdbJ0d724/edit
