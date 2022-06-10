import { TOKEN_KEY, USER_INFO } from "./config"
import { IUserInfo } from "./utils"

class LocalStorage {
    static set(key: string, value: string) {
        localStorage.setItem(key, value)
    }
    static get(key: string): string | null {
        return localStorage.getItem(key)
    }
    static remove(key: string) {
        localStorage.removeItem(key)
    }

    static getToken(): string {
        let token = this.get(TOKEN_KEY)
        if (token === undefined || token === null) {
            return ""
        }
        return token
    }
    static setToken(token: string) {
        this.set(TOKEN_KEY, token)
    }
    static removeToken() {
        this.remove(TOKEN_KEY)
    }

    static getUserInfo(): IUserInfo | null {
        let userInfo = this.get(USER_INFO)
        if (userInfo === undefined || userInfo === null) {
            return null
        }
        return JSON.parse(userInfo)
    }

    static setUserInfo(userInfo: IUserInfo) {
        this.set(USER_INFO, JSON.stringify(userInfo))
    }
    static removeUserInfo() {
        this.remove(USER_INFO)
    }
}

export default LocalStorage