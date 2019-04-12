package controllers

import (
    "encoding/json"
    "fmt"
    "github.com/astaxie/beego"
    "github.com/gorilla/websocket"
    "net/http"
)

/**
@author: agamgn
@date:	2019-04-12
*/


type ServersController struct {
    beego.Controller
}

func (s *ServersController)Get(){
    name:=s.GetString("name")
    if len(name)==0 {
        s.Redirect("/",302)
        return
    }
    s.Data["name"]=name
    s.TplName="Room.html"
}

type Client struct {
    conn *websocket.Conn	// 用户websocket连接
    name string
}

type Message struct {
    EventType byte	`json:"type"`		// 0表示用户发布消息；1表示用户进入；2表示用户退出
    Name string		`json:"name"`		// 用户名称
    Message string	`json:"message"`	// 消息
}

var (
    join = make(chan Client, 10)			// 用户加入通道
    leave = make(chan Client, 10)			// 用户退出通道
    message = make(chan Message, 10)		// 消息通道
    clients = make(map [Client] bool)		// 用户映射
)



func (s *ServersController)WsRoom(){
    name := s.GetString("name")
    if len(name) == 0 {
        s.Redirect("/", 302)
        return
    }
    // 检验http头中upgrader属性，若为websocket，则将http协议升级为websocket协议
    conn, err := (&websocket.Upgrader{}).Upgrade(s.Ctx.ResponseWriter, s.Ctx.Request, nil)
    if _, ok := err.(websocket.HandshakeError); ok {
        http.Error(s.Ctx.ResponseWriter, "Not a websocket handshake", 400)
        return
    } else if err != nil {
        return
    }
    //创建用户
    var client Client
    client.name=name
    client.conn=conn
    if !clients[client] {
        join <- client
    }
    // 当函数返回时，将该用户加入退出通道，并断开用户连接
    defer func() {
        leave <- client
        client.conn.Close()
    }()

    for {
        // 读取消息。如果连接断开，则会返回错误
        _, msgStr, err := client.conn.ReadMessage()
        // 如果返回错误，就退出循环
        if err != nil {
            break
        }
        //如果没有错误，则把用户发送的信息放入message通道中
        var msg Message
        msg.Name = client.name
        msg.EventType = 0
        msg.Message = string(msgStr)
        message <- msg
    }


}



func init() {
    go broad()
}
func broad() {
    for {
        // 哪个case可以执行，则转入到该case。都不可执行，则堵塞。
        select {
        case msg := <-message:
            str := fmt.Sprintf("broadcaster-----------%s send message: %s\n", msg.Name, msg.Message)
            fmt.Println(str)
            for client := range clients {
                data, err := json.Marshal(msg)
                if err != nil {
                    return
                }
                if client.conn.WriteMessage(websocket.TextMessage, data) != nil {
                }
            }
        // 有用户加入
        case client := <-join:
            str := fmt.Sprintf("broadcaster-----------%s join in the chat room\n", client.name)
            fmt.Println(str)
            clients[client] = true	// 将用户加入映射
            // 将用户加入消息放入消息通道
            var msg Message
            msg.Name = client.name
            msg.EventType = 1
            msg.Message = fmt.Sprintf("%s join in, there are %d preson in room", client.name, len(clients))
            message <- msg
        // 有用户退出
        case client := <-leave:
            str := fmt.Sprintf("broadcaster-----------%s leave the chat room\n", client.name)
            fmt.Println(str)
            if !clients[client] {
                break
            }
            delete(clients, client)
            // 将用户退出消息放入消息通道
            var msg Message
            msg.Name = client.name
            msg.EventType = 2
            msg.Message = fmt.Sprintf("%s leave, there are %d preson in room", client.name, len(clients))
            message <- msg
        }
    }
}
