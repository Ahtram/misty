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
[ok] 替 Misty 加上參數檢查邏輯，如果輸入參數錯誤會提醒使用者。
[ok] 替 Misty 設計 Config Sheet 以及 Config 結構，此結構用來設定這個機器人的工作內容以及身分。
[ok] 將不同身分的 Misty Bot 表單分開，並且透過參數在執行時傳入表單ID，去決定這個機器人的身分以及工作。
[ok] 格式化 help command 的顯示方式。讓使用者易於閱讀。
[ok] 支援 Beam 或者 hitbox (特定使用者頻道) 開始直播的 Discord 通知。
[ok] 使用 cid 命令來取得某個 Channel 的 ID。
[ok] config 現在支援指定 BroadcastChannelID ，設定機器人廣播用頻道。

config 新增變數 onlineNotify 設定上線時是否發送通知訊息。
broadcast功能，讓 misty 可以透過命令在管理者頻道內指定 misty 廣播到所有 Broadcast Channel.
限定特定頻道才能對 Misty 使用特定的管理者 Command. (update/broadcast) (權限)
支援 misty 開啟時設定語系。
支援 misty 在執行時期可以接受命令更換語系。
測試 golong timer 以及隨機自言自語功能的可能性。

讀取TET的資料表單，提供玩家查詢遊戲內的合成公式或者道具資訊(尚未確定)。

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

//讓 uni 加入自己 server 用的連結
https://discordapp.com/oauth2/authorize?client_id=240432374574350338&scope=bot&permissions=0

//執行 Misty
./misty -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.Cus7Sw.n-suc_aXypKw-EnkRw8kA3TMU4Y" -c 1H46UgwUKfg7OqE31uxj5ko_B_a_E7Y600W8eM4V2EGI
or
./misty -e [email] -p [password]
or 
./misty -e [email] -p [password] -t "Bot MjMxMTA1MTQ4MDc0NzIxMjgw.Cus7Sw.n-suc_aXypKw-EnkRw8kA3TMU4Y"

//執行Uni
./misty -t "Bot MjQwNDMyMzc0NTc0MzUwMzM4.CvDPpg.qqrJHoy27PvltMjsbTIFsdg9QpQ" -c 1fBpU4E9vE6BQVYX8ho1GW6NGm-Exg9rzKo4nMWlbcck

使用 tmux 指令來背景常駐運作 misty
tmux 教學: https://gist.github.com/MohamedAlaa/2961058

擴充命令用 Google Sheet
https://docs.google.com/spreadsheets/d/1haLbQuE7TtF79_J2XLbzFRYbAkfGRCmrXxwdbJ0d724/edit
