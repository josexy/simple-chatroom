import { MsgType } from "./websocket"

export interface IResponseData<T> {
    code: number,
    msg: string,
    data?: T
}

export interface IWsResponse {
    status: number
    msg?: string
    data?: IMessage
    list?: Array<IClientInfo>
}

export interface IRoomInfo {
    id: number
    title: string
    description: string
    online_clients: number
}

export interface IUserInfo {
    id: number,
    username: string,
    nickname: string,
    avatar_id: number,
}

export interface IUserTokenInfo {
    token: string,
    expire: Date
}

// 在线用户信息
export interface IClientInfo {
    uid: number
    username: string
    nickname: string
    room_id: number
    avatar_id: number
}

export interface IMessage {
    uid?: number
    username?: string
    nickname?: string
    room_id?: number
    avatar_id?: number
    to_user_id?: number
    content?: string
    image_url?: string
}

export interface IHistoryMessage {
    msg_type?: MsgType
    id?: number,
    uid?: number,
    to_user_id?: number,
    room_id?: number,
    avatar_id?: number
    content?: string,
    image_url?: string
    username?: string
    nickname?: string
}

export interface IOnlineClientCount {
    count: number
}

export interface IUploadImage {
    image_url: string
}

const copyToClipboard = (text: string) => {
    try {
        if (navigator.clipboard && window.isSecureContext) {
            (async () => {
                await navigator.clipboard.writeText(text)
            })()
        } else {
            const textarea = document.createElement('textarea')
            textarea.value = text
            document.body.appendChild(textarea)
            textarea.select()
            document.execCommand('copy')
            document.body.removeChild(textarea)
        }
    } catch (error) {
        console.log(error);
    }
}

export { copyToClipboard }
