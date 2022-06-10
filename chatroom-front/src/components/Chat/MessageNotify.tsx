import { useEffect, useState } from "react"

interface IProps {
    nickname?: string
    onlineOrOffline: boolean
}

export default function MessageNotify({ nickname, onlineOrOffline }: IProps) {

    const [time, setTime] = useState('')

    useEffect(() => {
        setTime(new Date(Date.now()).toLocaleString())
    }, [])

    return (
        <div className="text-center" style={{ "backgroundColor": "lightgray" }}>
            {onlineOrOffline ? <span style={{ "color": "green" }}>User Coming</span>
                : <span style={{ "color": "red" }}>User Left</span>
            }
            :<span className="fw-bolder ms-2">{`[ ${nickname} ]`}</span>
            <span className="ms-2" style={{ "fontSize": "10px" }}>
                {time}
            </span>
        </div>
    )
}