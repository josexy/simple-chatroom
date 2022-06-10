import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios'
import { BASE_URL, TIME_OUT } from './config'
import LocalStorage from './local'
import { IUserTokenInfo } from './utils'

interface IReqCallback {
    callback: (token: string) => void
}

class AxiosRequest {

    public instance: AxiosInstance
    private isRefreshing: boolean = false
    private peedingRequest: Array<IReqCallback> = []

    constructor(config: AxiosRequestConfig) {
        this.instance = axios.create(config)
        this.instance.interceptors.request.use(
            config => {
                const token = LocalStorage.getToken()
                if (token !== "" && config.headers) {
                    config.headers['Authorization'] = "Bearer " + token
                }
                return config
            },
            err => {
                return err
            }
        )
        this.instance.interceptors.response.use(
            response => {
                return response
            },
            error => {
                if (!error.response) {
                    return Promise.reject(error)
                }
                if (error.response.status === 401 && !error.config.url.includes('/auth/refresh_token')) {
                    const { config } = error.response
                    if (!this.isRefreshing) {
                        this.isRefreshing = true
                        return this.refreshToken().then(resp => {
                            if (resp.data.code === 1000 && resp.data.data) {
                                const token = resp.data.data.token
                                LocalStorage.setToken(token)

                                if (config.headers) {
                                    config.headers['Authorization'] = `Bearer ${token}`
                                }

                                this.peedingRequest.forEach((req) => {
                                    req.callback(token)
                                })
                                this.peedingRequest = []
                                return this.instance(config)
                            }
                        }).catch(err => {
                            LocalStorage.removeToken()
                            LocalStorage.removeUserInfo()
                            setTimeout(() => {
                                window.location.replace('/')
                            }, 2500);
                            return Promise.reject(err)
                        }).finally(() => {
                            this.isRefreshing = false
                        })
                    } else {
                        return new Promise(resolve => {
                            this.peedingRequest.push({
                                callback: (token: string) => {
                                    if (config.headers)
                                        config.headers['Authorization'] = `Bearer ${token}`
                                    resolve(this.instance(config))
                                }
                            })
                        })
                    }
                } else if (error.response && error.response.status === 403) {
                }
                return Promise.reject(error)
            }
        )
    }

    refreshToken() {
        interface IResponseData {
            code: number,
            msg: string,
            data?: IUserTokenInfo
        }

        return this.get<IResponseData>("/auth/refresh_token", {
            headers: { "Authorization": `Bearer ${LocalStorage.getToken()}` }
        })
    }

    async request<T>(config: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.instance.request<any, any>(config)
    }

    get<T>(url: string, config?: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.request<T>({ url: url, ...config, method: 'GET' })
    }

    post<T>(url: string, config?: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.request<T>({ url: url, ...config, method: 'POST' })
    }

    delete<T>(url: string, config?: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.request<T>({ url: url, ...config, method: 'DELETE' })
    }

    put<T>(url: string, config?: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.request<T>({ url: url, ...config, method: "PUT" })
    }

    patch<T>(url: string, config?: AxiosRequestConfig<any>): Promise<AxiosResponse<T>> {
        return this.request<T>({ url: url, ...config, method: 'PATCH' })
    }
}

const request = new AxiosRequest({
    baseURL: BASE_URL,
    timeout: TIME_OUT,
    headers: {
        "Content-Type": "application/json",
    },
})

export default request

