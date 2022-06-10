import { useEffect, useState } from "react";
import { Badge, Card } from "react-bootstrap";
import { useNavigate } from "react-router-dom";
import LocalStorage from "../api/local";
import request from "../api/request";
import { IResponseData, IRoomInfo } from "../api/utils";
import { GPanel282 } from "./GPanel";

export default function GHome() {

    const [rooms, setRooms] = useState<IRoomInfo[]>([])
    const navigate = useNavigate()

    useEffect(() => {

        // 是否登录
        if (!LocalStorage.getUserInfo()) {
            navigate('/')
            return
        }

        request.get<IResponseData<IRoomInfo[]>>("/auth/rooms").then(res => {
            if (res.data.code === 1000) {
                if (res.data.data) {
                    setRooms(res.data.data)
                }
            }
        }).catch(err => {
            console.log(err)
        })
    }, [])

    return (
        <GPanel282>
            {
                rooms && rooms.map((val, index) =>
                    <div className="d-inline-flex m-2" key={index}>
                        <Card
                            bg={'light'}
                            key={index}
                            style={{ width: '18rem' }}
                            className="mb-2"
                        >
                            <Card.Header>Room {val.id}</Card.Header>
                            <Card.Body>
                                <Card.Title>Room {val.title} </Card.Title>
                                <div style={{ "color": "gray" }}>
                                    {val.description}
                                </div>
                                <div>
                                    Online Users: <Badge bg={"primary"}>{val.online_clients}</Badge>
                                </div>
                                <Card.Link href={`/room/${val.id}`}>Enter Room</Card.Link>
                            </Card.Body>
                        </Card>
                    </div>
                )
            }

        </GPanel282 >
    )
}