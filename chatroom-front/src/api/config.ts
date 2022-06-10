
let BASE_URL = ''

const TOKEN_KEY = "jwt_token"
const USER_INFO = "_userinfo"

const TIME_OUT = 5000

const PRODUCTION_URL = "/api/v1"
const DEPLOYMENT_URL = "http://127.0.0.1:10086/api/v1"

const PRODUCTION_WS_URL = `ws://${document.location.host}/api/v1/ws`
const DEPLOYMENT_WS_URL = "ws://127.0.0.1/api/v1/ws"
let WS_URL = ''

if (process.env.NODE_ENV === 'development') {
    BASE_URL = DEPLOYMENT_URL
    WS_URL = DEPLOYMENT_WS_URL
} else if (process.env.NODE_ENV === 'production') {
    BASE_URL = PRODUCTION_URL
    WS_URL = PRODUCTION_WS_URL
} else {
    BASE_URL = DEPLOYMENT_URL
}

export { BASE_URL, TIME_OUT, TOKEN_KEY, USER_INFO, WS_URL }
