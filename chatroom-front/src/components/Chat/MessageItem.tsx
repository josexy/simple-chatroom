
interface IProps {
    from_me: boolean
    nickname?: string
    content?: string
    avatar_id?: number
    image_url?: string
}

export default function MessageItem({ avatar_id, nickname, from_me, content, image_url }: IProps) {
    return (
        <div>
            <div className={from_me ? 'message-from-me-name' : 'message-from-other-name'}>
                <img
                    className="message-user-image"
                    alt=""
                    src={`/qqimg/${avatar_id ? avatar_id : 1}.jpg`}
                />
                <span className="fw-bold ms-2 me-2 align-self-center">
                    {nickname}
                </span>
            </div>
            <div className={`message-item ${from_me ? 'message-from-me' : 'message-from-other'}`}>
                <div>
                    <div>
                        {image_url && <img className="message-content-image" alt="" src={image_url} />}
                    </div>
                    <div>
                        {content ? content : ""}
                    </div>
                </div>
            </div>
        </div>
    )
}