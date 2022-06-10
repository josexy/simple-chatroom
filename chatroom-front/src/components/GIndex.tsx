import { Alert, Button, Form, InputGroup } from "react-bootstrap";
import { GPanel444 } from "./GPanel";
import { useEffect, useState } from "react";
import request from "../api/request";
import { IOnlineClientCount, IResponseData, IUserInfo, IUserTokenInfo } from "../api/utils";
import { useNavigate } from "react-router-dom";
import LocalStorage from "../api/local";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBookOpen, faLock, faSignInAlt, faUser, faUserEdit } from "@fortawesome/free-solid-svg-icons";

function GIndex() {

    const navigate = useNavigate()
    const [onlineClientCount, setOnlineClientCount] = useState(0)

    const [username, setUsername] = useState('')
    const [password, setPassword] = useState('')
    const [nickname, setNickname] = useState('')

    useEffect(() => {
        // 如果已经登录，则跳转到房间界面
        if (LocalStorage.getUserInfo() !== null) {
            navigate("/rooms")
        } else {
            // 获取在线用户总数
            request.get<IResponseData<IOnlineClientCount>>("/online").then(res => {
                if (res.data.code === 1000) {
                    if (res.data.data) {
                        setOnlineClientCount(res.data.data.count)
                    }
                }
            }).catch(err => {
                console.log(err);
            })
        }
    }, [])

    const login = () => {
        let u = username.trim()
        let p = password.trim()
        if (u.length < 3 || u.length > 20 || p.length < 3 || p.length > 20) {
            alert('username or password invalid!')
            return
        }

        // 登录
        request.post<IResponseData<IUserTokenInfo>>("/user/login", {
            data: {
                username: u,
                password: p,
            }
        }).then(res => {
            switch (res.data.code) {
                case 1000:
                    if (res.data.data) {
                        LocalStorage.setToken(res.data.data.token)
                        // 获取用户信息
                        request.get<IResponseData<IUserInfo>>("/auth/user").then(res => {
                            if (res.data.code === 1000) {
                                if (res.data.data) {
                                    LocalStorage.setUserInfo(res.data.data)
                                }
                            }
                        }).catch(err => {
                            console.log(err)
                        })
                    }
                    // 登录成功后进入选择房间界面
                    setTimeout(() => {
                        navigate("/rooms")
                    }, 1000)
                    break
                default:
                    alert('login failed!')
            }
        }).catch(err => {
            alert('login failed!')
            console.log(err);
        })
    }

    const register = () => {
        let u = username.trim()
        let p = password.trim()
        let n = nickname.trim()
        if (u.length < 3 || u.length > 20 || n.length < 3 || n.length > 20 || p.length < 3 || p.length > 20) {
            alert('username or password or nickname invalid!')
            return
        }
        request.post<IResponseData<IUserInfo>>("/user/register", {
            data: {
                username: username,
                nickname: nickname,
                password: password,
                avatar_id: 1, // TODO
            }
        }).then(res => {
            switch (res.data.code) {
                case 1000:
                    alert('register successfully!')
                    break
                default:
                    alert('register failed!')
                    break
            }
        }).catch(err => {
            console.log(err);
        })
    }

    return (
        <GPanel444>
            <Form>
                <Form.Group className="mb-2">
                    <Form.Label>Username</Form.Label>
                    <InputGroup>
                        <InputGroup.Text>
                            <FontAwesomeIcon icon={faUser} />
                        </InputGroup.Text>
                        <Form.Control minLength={3} maxLength={20}
                            value={username}
                            onChange={(e: any) => setUsername(e.target.value)}
                            type="text"
                            placeholder="username characters between 3 and 20" />
                    </InputGroup>
                </Form.Group>

                <Form.Group className="mb-2">
                    <Form.Label>Password</Form.Label>
                    <InputGroup>
                        <InputGroup.Text>
                            <FontAwesomeIcon icon={faLock} />
                        </InputGroup.Text>
                        <Form.Control minLength={3} maxLength={20}
                            value={password}
                            onChange={(e: any) => setPassword(e.target.value)}
                            type="password"
                            placeholder="password characters between 3 and 20" />
                    </InputGroup>
                </Form.Group>

                <Form.Group className="mb-2">
                    <Form.Label>{`Nickname(for register)`}</Form.Label>
                    <InputGroup>
                        <InputGroup.Text>
                            <FontAwesomeIcon icon={faBookOpen} />
                        </InputGroup.Text>
                        <Form.Control minLength={3} maxLength={20}
                            value={nickname}
                            onChange={(e: any) => setNickname(e.target.value)}
                            type="text"
                            placeholder="nickname characters between 3 and 20" />
                    </InputGroup>
                </Form.Group>

                <div className={"d-grid gap-2 mb-2"}>
                    <Button variant="primary" type="button" onClick={register}>
                        Register
                        <FontAwesomeIcon className="ms-2" icon={faUserEdit} />
                    </Button>
                    <Button variant="success" type="button" onClick={login}>
                        Login
                        <FontAwesomeIcon className="ms-2" icon={faSignInAlt} />
                    </Button>
                </div>
                <Alert variant="primary">
                    Current online users : <span
                        style={{ "color": "red" }} className="fw-bolder">{onlineClientCount}
                    </span>
                </Alert>
            </Form>
        </GPanel444>
    )
}

export default GIndex