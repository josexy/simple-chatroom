import { useEffect, useRef, useState } from "react"
import { Row } from "react-bootstrap"
import LocalStorage from "../../api/local"
import { IHistoryMessage, IUserInfo } from "../../api/utils"
import { MsgType } from "../../api/websocket"
import MessageItem from "./MessageItem"
import MessageNotify from "./MessageNotify"

interface IProps {
    messages: Array<IHistoryMessage>
}

function MessageServerClosed() {
    return (
        <div className="h5 fw-bolder text-center"
            style={{ color: "red", backgroundColor: "lightgrey" }}>
            Server Closed!
        </div>
    )
}

export default function MessageList({ messages }: IProps) {

    const ref = useRef<HTMLDivElement>(null)
    const [userInfo, setUserInfo] = useState<IUserInfo>()

    useEffect(() => {
        const userinfo = LocalStorage.getUserInfo()
        if (userinfo)
            setUserInfo(userinfo)
        scrollToBottom()
    }, [])

    useEffect(() => {
        scrollToBottom()
    }, [messages])

    // 滚动到底部
    const scrollToBottom = () => {
        if (ref.current) {
            ref.current.scrollTop = ref.current.scrollHeight
        }
    }

    return (
        <Row>
            <div ref={ref} className='message-list'>
                {
                    messages.map((message, index) => {
                        if (message.msg_type === MsgType.ServerPanic)
                            return (
                                <MessageServerClosed key={index} />
                            )
                        else if (message.msg_type === MsgType.Online || message.msg_type === MsgType.Offline)
                            return (<MessageNotify
                                key={index}
                                nickname={message.nickname}
                                onlineOrOffline={message.msg_type === MsgType.Online}
                            />)
                        else return (
                            <MessageItem
                                key={index}
                                avatar_id={message.avatar_id}
                                nickname={message.nickname}
                                from_me={message.uid === userInfo?.id}
                                content={message.content}
                                image_url={message.image_url}
                            />
                        )
                    })
                }
            </div>
        </Row>
    )
}