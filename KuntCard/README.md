## 巫师3昆塔牌 web版本

1. 用户第一次使用websocket链接到服务器,发送进入房间的消息给服务器
2. 服务器为用户分配房间，先查找是否有处于等待的房间，有的话就分配过去，此房间就可以开始PK了,设置isRunning = true并且isWait = false；如果没有的话，新开一个房间，把用户加进去后，房间处于等待状态，等待另一个用户加入后开始PK
3. 第二个用户加入后，开始PK。服务器分别向该房间中的两个用户返回各自的牌组。
