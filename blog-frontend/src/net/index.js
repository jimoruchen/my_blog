import axios from "axios";
import {ElMessage} from "element-plus";
import router from "@/router";

const authItemName = "authorize"

export const accessHeader = () => {
    return {
        'Authorization': `Bearer ${takeAccessToken()}`
    }
}

const defaultError = (error) => {
    console.error(error)
    const status = error.response.status
    if (status === 429) {
        ElMessage.error(error.response.data.message)
    } else {
        ElMessage.error('发生了一些错误，请联系管理员')
    }
}

const defaultFailure = (message, status, url) => {
    console.warn(`请求地址: ${url}, 状态码: ${status}, 错误信息: ${message}`)
    ElMessage.warning(message)
}

function takeAccessToken() {
    const str = localStorage.getItem(authItemName) || sessionStorage.getItem(authItemName);
    if(!str) return null
    try {
        const authObj = JSON.parse(str)
        if(new Date(authObj.expire) <= new Date()) {
            deleteAccessToken()
            ElMessage.warning("登录状态已过期，请重新登录！")
            return null
        }
        return authObj.token
    } catch (error) {
        console.error('解析token时出错:', error);
        deleteAccessToken()
        ElMessage.warning("登录状态异常，请重新登录！")
        return null
    }
}

function storeAccessToken(remember, token, expire){
    const authObj = {
        token: token,
        expire: expire
    }
    const str = JSON.stringify(authObj)
    if(remember)
        localStorage.setItem(authItemName, str)
    else
        sessionStorage.setItem(authItemName, str)
}

function deleteAccessToken(redirect = false) {
    localStorage.removeItem(authItemName)
    sessionStorage.removeItem(authItemName)
    if(redirect) {
        router.push({ name: 'welcome-login' })
    }
}

function internalPost(url, data, headers, success, failure, error = defaultError){
    axios.post(url, data, { headers: headers, withCredentials: true }).then(({data}) => {
        if(data.code === 1000) {
            success(data.data)
        } else if(data.code === 401) {
            failure('登录状态已过期，请重新登录！')
            deleteAccessToken(true)
        } else {
            failure(data.message, data.code, url)
        }
    }).catch(err => error(err))
}

function internalGet(url, headers, success, failure, error = defaultError){
    axios.get(url, { headers: headers, withCredentials: true }).then(({data}) => {
        if(data.code === 1000) {
            success(data.data)
        } else if(data.code === 401) {
            failure('登录状态已过期，请重新登录！')
            deleteAccessToken(true)
        } else {
            failure(data.message, data.code, url)
        }
    }).catch(err => error(err))
}

function internalDelete(url, headers, success, failure, error = defaultError){
    axios.delete(url, { headers: headers, withCredentials: true }).then(({data}) => {
        if(data.code === 1000) {
            success(data.data)
        } else if(data.code === 401) {
            failure('登录状态已过期，请重新登录！')
            deleteAccessToken(true)
        } else {
            failure(data.message, data.code, url)
        }
    }).catch(err => error(err))
}

function login(username, password, remember, success, failure = defaultFailure){
    internalPost('/api/auth/login', {
        username: username,
        password: password
    }, {
        'Content-Type': 'application/x-www-form-urlencoded'
    }, (data) => {
        storeAccessToken(remember, data.token, data.expire)
        ElMessage.success(`登录成功，欢迎 ${data.username} 来到我们的系统`)
        success(data)
    }, failure)
}

function post(url, data, success, failure = defaultFailure) {
    internalPost(url, data, accessHeader() , success, failure)
}

function logout(success, failure = defaultFailure){
    get('/api/auth/logout', () => {
        deleteAccessToken()
        ElMessage.success(`退出登录成功，欢迎您再次使用`)
        success()
    }, failure)
}

function get(url, success, failure = defaultFailure, needAuth = true) {
    const headers = needAuth ? accessHeader() : {};
    internalGet(url, headers, success, failure);
}
function deleteRequest(url, success, failure = defaultFailure) {
    internalDelete(url, accessHeader(), success, failure)
}

function unauthorized() {
    return !takeAccessToken()
}

export { post, get, login, logout, unauthorized, deleteRequest }
