import { Badge } from "react-bootstrap";
import { IClientInfo } from "../../api/utils";

interface IProps {
    uid?: number
    unreadMsgCount?: number
    clientInfo: IClientInfo
}

export default function Messager({ uid, unreadMsgCount, clientInfo }: IProps) {
    return (
        <div style={{ "cursor": "pointer" }}>
            <div className="d-flex align-items-center">
                <img
                    className="message-user-image"
                    alt=""
                    src={`/qqimg/${clientInfo.avatar_id}.jpg`}
                />
                <span
                    className="fw-bold">
                    {clientInfo.nickname}
                </span>
                {uid !== undefined && uid === clientInfo.uid &&
                    <Badge className="ms-1" bg={"danger"}>Me</Badge>}
                {unreadMsgCount !== undefined && unreadMsgCount > 0 &&
                    <Badge className="ms-auto" bg={"success"}>{unreadMsgCount}</Badge>}
            </div>
        </div>
    );
}