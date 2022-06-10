import { faImage } from "@fortawesome/free-solid-svg-icons";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { ChangeEvent, useRef, useState } from "react"
import { Row } from "react-bootstrap"
import InputEmoji from "react-input-emoji";
import request from "../../api/request";
import { IResponseData, IUploadImage } from "../../api/utils";

interface IProps {
    onSend: any
}

export default function MessageInputBox({ onSend }: IProps) {

    const [messageText, setMessageText] = useState('')
    const [imageUrl, setImageUrl] = useState('')
    const fileRef = useRef<HTMLInputElement>(null)

    const onEnter = () => {
        onSend(messageText, imageUrl)
        setMessageText('')
        setImageUrl('')
    }

    const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
        if (e.target.files) {
            if (e.target.files.length > 0) {
                const file = e.target.files.item(0)
                if (file) {
                    if (file.type.indexOf("image") !== -1) {
                        const data = new FormData()
                        data.append('uploadImage', file)

                        request.post<IResponseData<IUploadImage>>('/auth/img_upload', { data: data })
                            .then(res => {
                                if (res.data.code === 1000 && res.data.data) {
                                    setImageUrl(res.data.data.image_url)
                                }
                            }).catch(err => {
                                console.log(err);
                            })
                    } else {
                        alert('the file is not an image!')
                    }
                }
            }
        }
        if (fileRef.current)
            fileRef.current.value = ''
    }

    const uploadImage = () => {
        if (fileRef.current)
            fileRef.current.click()
    }

    return (
        <Row className="mt-2 message-input-box">
            <div className="ms-3 me-3 mt-2">
                <div className="upload-image" onClick={uploadImage}>
                    <FontAwesomeIcon size="2x" icon={faImage} />
                </div>
            </div>
            <div className="mb-2">
                <InputEmoji
                    value={messageText}
                    onChange={(text: any) => setMessageText(text)}
                    cleanOnEnter
                    onEnter={onEnter}
                    theme="light"
                    height={100}
                    borderRadius={5}
                    placeholder="Enter a friendly message :)"
                />
            </div>
            <label htmlFor="uploadImage"></label>
            <input
                id="uploadImage"
                ref={fileRef}
                type={"file"}
                onChange={handleChange}
                hidden
            />
        </Row >
    )
}