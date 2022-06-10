import { WS_URL } from "./config"
import LocalStorage from "./local"
import { IUserInfo } from "./utils"

enum MsgType {
    Heartbeat,
    Online,
    Offline,
    Send,
    OnlineClients,
    PrivateChat,
    ServerPanic,
}

interface WsCallback {
    onOpen: (event: Event) => any
    onClose: (event: CloseEvent) => any
    onMessage: (event: MessageEvent) => any
    onError: (event: Event) => any
}

class WsClient {

    public wsConn: WebSocket
    public intervalHeartbeat: any
    public intervalOnlineClients: any
    public userInfo: IUserInfo | null

    constructor(callback: WsCallback) {
        this.userInfo = LocalStorage.getUserInfo()
        this.wsConn = new WebSocket(WS_URL)
        this.wsConn.onopen = callback.onOpen
        this.wsConn.onclose = callback.onClose
        this.wsConn.onerror = callback.onError
        this.wsConn.onmessage = callback.onMessage
    }

    isClosed() {
        /*
            0 (WebSocket.CONNECTING) 正在链接中
            1 (WebSocket.OPEN) 已经链接并且可以通讯
            2 (WebSocket.CLOSING) 连接正在关闭
            3 (WebSocket.CLOSED) 连接已关闭或者没有链接成功
         */
        return this.wsConn.readyState === 2 || this.wsConn.readyState === 3
    }

    close() {
        this.stopHeartbeat()
        this.stopOnlineClients()
        this.wsConn.close()
    }

    stopHeartbeat() {
        clearInterval(this.intervalHeartbeat)
    }

    startHeartbeat() {
        this.intervalHeartbeat = setInterval(() => {
            this.send({
                status: MsgType.Heartbeat,
                msg: "heartbeat"
            })
        }, 5000)
    }

    stopOnlineClients() {
        clearInterval(this.intervalOnlineClients)
    }

    refreshOnlineClients(roomId?: string) {
        this.sendMessage(MsgType.OnlineClients, roomId)
        this.intervalOnlineClients = setInterval(() => {
            this.sendMessage(MsgType.OnlineClients, roomId)
        }, 3000)
    }

    addRoom(roomId?: string) {
        this.sendMessage(MsgType.Online, roomId)
    }

    sendMessage(status: MsgType, roomId?: string, content?: string, to_user_id?: number, image_url?: string) {
        if (this.userInfo === null) return

        const uid = this.userInfo.id
        const room_id = Number.parseInt(roomId !== undefined ? roomId : "0")
        const username = this.userInfo.username
        const nickname = this.userInfo.nickname
        const avatar_id = this.userInfo.avatar_id

        this.send({
            status: status,
            data: {
                uid: uid,
                room_id: room_id,
                username: username,
                nickname: nickname,
                avatar_id: avatar_id,
                content: content,
                to_user_id: to_user_id,
                image_url: image_url,
            }
        })
    }
    send(data: any) {
        this.wsConn.send(JSON.stringify(data))
    }
}

export { WsClient, MsgType }