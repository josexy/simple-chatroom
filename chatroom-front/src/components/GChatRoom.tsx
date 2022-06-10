import { useEffect, useRef, useState } from "react";
import { Badge, Button, Card, Col, ListGroup, Row, Stack } from "react-bootstrap";
import { useNavigate, useParams } from "react-router-dom";
import request from "../api/request";
import { IClientInfo, IHistoryMessage, IMessage, IResponseData, IWsResponse } from "../api/utils";
import MessageList from "./Chat/MessageList";
import MessageInputBox from "./Chat/MessageInputBox";
import { MsgType, WsClient } from "../api/websocket";
import Messager from "./Chat/Messager";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBackward, faUserGroup } from "@fortawesome/free-solid-svg-icons";
import LocalStorage from "../api/local";

export default function GChatRoom() {

    const navigate = useNavigate()
    const params = useParams()
    const [messages, setMessages] = useState<IHistoryMessage[]>([])
    const [onlineClients, setOnlineClients] = useState<IClientInfo[]>([])
    const [itemIndex, setItemIndex] = useState(-1)
    const [uid, setUid] = useState<number>()

    // https://zh-hans.reactjs.org/docs/hooks-faq.html#why-am-i-seeing-stale-props-or-state-inside-my-function
    const toClient = useRef<IClientInfo | null>()
    const unreadMsgCount = useRef<Map<number | undefined, number | undefined>>(new Map())

    const ws = useRef<WsClient>();

    const onWsOpen = (e: Event) => {
        console.log("打开websocket");

        if (!ws.current) return

        // 建立websocket连接后，加入到对应的房间中
        ws.current.addRoom(params.id)
        // 定时发送心跳包
        ws.current.startHeartbeat()
        // 定时获取当前房间内的在线用户
        ws.current.refreshOnlineClients(params.id)
    }

    const onWsClose = (e: CloseEvent) => {
        console.log("关闭websocket");

        if (ws.current) {
            ws.current.close()
        }
    }

    const onWsError = (e: Event) => {
        console.log("错误websocket");
        addNewMessage(MsgType.ServerPanic, {})
    }

    const onWsMessage = (e: MessageEvent) => {
        const data: IWsResponse = JSON.parse(e.data)
        if (data.status === MsgType.Heartbeat) {
            // 心跳包
        } else if (data.status === MsgType.OnlineClients) {
            // 获取当前房间内在线用户
            if (data.list) {
                setOnlineClients(data.list)
            }
        } else if (data.status === MsgType.Online) {
            // 新用户进入房间
            console.log('新用户进入房间');
            if (data.data) {
                if (unreadMsgCount.current && data.data.uid !== ws.current?.userInfo?.id) {
                    unreadMsgCount.current.set(data.data.uid, 0)
                }
                if (!toClient.current) {
                    addNewMessage(MsgType.Online, data.data)
                }
            }
        } else if (data.status === MsgType.Send) {
            // 新消息到达
            console.log('新消息');
            if (data.data) {
                // 是否切换到群聊模式
                if (!toClient.current)
                    addNewMessage(MsgType.Send, data.data)
            }
        } else if (data.status === MsgType.Offline) {
            // 用户下线
            console.log('用户下线');
            if (data.data) {
                if (unreadMsgCount.current) {
                    unreadMsgCount.current.delete(data.data.uid)
                }
                if (toClient.current && toClient.current.uid === data.data.uid) {
                    exitRoom()
                    return
                }
                if (!toClient.current)
                    addNewMessage(MsgType.Offline, data.data)
            }
        } else if (data.status === MsgType.PrivateChat) {
            // 私聊
            console.log('私聊')
            if (data.data) {
                if (toClient.current &&
                    ((ws.current?.userInfo?.id === data.data.uid) ||
                        (toClient.current.uid === data.data.uid &&
                            ws.current?.userInfo?.id === data.data.to_user_id))) {
                    addNewMessage(MsgType.PrivateChat, data.data)
                }
                if (unreadMsgCount.current) {
                    if (data.data.uid !== ws.current?.userInfo?.id) {
                        if (toClient.current && toClient.current.uid === data.data.uid) {
                            return
                        }
                        const count = unreadMsgCount.current.get(data.data.uid)
                        if (count !== undefined) {
                            unreadMsgCount.current.set(data.data.uid, count + 1)
                        }
                    }
                }
            }
        }
    }

    const addNewMessage = (typ: MsgType, data: IMessage) => {
        const msg: IHistoryMessage = {
            msg_type: typ,
            uid: data.uid,
            username: data.username,
            nickname: data.nickname,
            avatar_id: data.avatar_id,
            content: data.content,
            image_url: data.image_url,
        }

        setMessages(oldMessage => [...oldMessage, msg])
    }

    const getHistoryMessages = (url: string, params?: any) => {
        // 获取历史聊天记录
        request.get<IResponseData<IHistoryMessage[]>>(url, { params: params }).then(res => {
            console.log(res.data)
            if (res.data.code === 1000) {
                if (res.data.data) {
                    setMessages(res.data.data)
                }
            }
        }).catch(err => {
            console.log(err)
        })
    }

    useEffect(() => {

        if (!LocalStorage.getUserInfo()) {
            console.log('not login');
            navigate('/')
            return
        }

        // websocket连接
        ws.current = new WsClient({
            onOpen: onWsOpen,
            onClose: onWsClose,
            onError: onWsError,
            onMessage: onWsMessage,
        })

        getHistoryMessages(`/auth/room/${params.id}`)

        setUid(ws.current.userInfo?.id)

        return () => {
            if (ws.current)
                ws.current.close()
        }

    }, [params.id])

    const onSendMessage = (message: string, imageUrl: string) => {
        if (message.replace(/(^\s*)|(\s*$)/g, "").length === 0) {
            if (imageUrl.length === 0) {
                return
            }
        }

        if (ws.current && !ws.current.isClosed()) {
            // 群聊还是私聊
            if (toClient.current) {
                ws.current.sendMessage(MsgType.PrivateChat, params.id, message, toClient.current.uid, imageUrl)
            } else {
                ws.current.sendMessage(MsgType.Send, params.id, message, 0, imageUrl)
            }
        }
    }

    const exitRoom = () => {
        console.log("exit room: ", params.id, toClient)
        if (toClient.current) {
            setItemIndex(-1)
            getHistoryMessages(`/auth/room/${params.id}`)

            toClient.current = null
            return
        }
        if (ws.current) {
            ws.current.close()
            navigate('/rooms')
        }
    }

    const messagerChange = (index: number, clientInfo: IClientInfo) => {
        if (!ws.current) return
        if (clientInfo.uid === ws.current.userInfo?.id) return
        if (index === itemIndex) return

        if (unreadMsgCount.current) {
            unreadMsgCount.current.set(clientInfo.uid, 0)
        }
        setItemIndex(index)
        getHistoryMessages(`/auth/room/p/${params.id}/${clientInfo.uid}`, { uid: ws.current.userInfo?.id })
        toClient.current = clientInfo
    }

    const logout = () => {
        console.log('logout');
        request.delete('/auth/user').then(res => {
            LocalStorage.removeToken()
            LocalStorage.removeUserInfo()
            navigate("/")
        }).catch(err => {
            console.log(err);
        })
    }

    return (
        <Row>
            <Col md={3} lg={3}>
                <div className="overflow-auto" style={{ maxHeight: "60vh" }}>
                    <ListGroup className="mb-2">
                        <ListGroup.Item className="text-center">
                            <FontAwesomeIcon className="me-1" icon={faUserGroup} />
                            <span className="me-1">Online Users</span>
                            <Badge bg="primary">
                                {onlineClients ? onlineClients.length : 0}
                            </Badge>
                        </ListGroup.Item>
                        {
                            onlineClients && onlineClients.map((val, index) =>
                                <ListGroup.Item
                                    key={index}
                                    as={"li"}
                                    action
                                    active={itemIndex === index}
                                    onClick={(e: any) => messagerChange(index, val)}
                                >
                                    <Messager
                                        uid={uid}
                                        clientInfo={val}
                                        unreadMsgCount={unreadMsgCount.current && unreadMsgCount.current.get(val.uid)}
                                    />
                                </ListGroup.Item>)
                        }
                    </ListGroup>
                </div>
            </Col>
            <Col md={9} lg={9}>
                <Card>
                    <Card.Header>
                        <div className="d-flex">
                            <Button
                                variant="outline-danger"
                                size={"sm"}
                                onClick={exitRoom}
                            >
                                <FontAwesomeIcon className="me-1" icon={faBackward} />
                                Back
                            </Button>
                            <div className="fw-bold align-self-center ms-auto">
                                {
                                    toClient.current ? `User: ${toClient.current.nickname}` :
                                        <span>
                                            {"Room " + params.id}
                                        </span>
                                }
                            </div>
                            <div className="ms-auto">
                                <Button
                                    variant="link"
                                    size={"sm"}
                                    onClick={logout}
                                >
                                    Logout
                                </Button>
                            </div>
                        </div>
                    </Card.Header>
                    <Card.Body>
                        <MessageList messages={messages} />
                        <MessageInputBox onSend={onSendMessage} />
                    </Card.Body>
                </Card>
            </Col>
        </Row>
    )
}
