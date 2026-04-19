/**
 * 文件: request-B-CfM8d9.js → request.js
 * 功能: HTTP 请求封装
 *
 * ⚠️ 这是简化版本的请求封装，基于源码中的使用方式重构
 * 原始文件 (132KB) 包含大量 Vue 组件库代码，无法完整反混淆
 */

const TOKEN_COOKIE_NAME = "webvpn-token";

/**
 * 从 Cookie 获取值
 */
function getCookie(name) {
    const match = document.cookie.match(new RegExp(`(^| )${name}=([^;]+)`));
    return match ? match[2] : null;
}

/**
 * 核心请求函数
 */
async function request(options) {
    const {
        url,
        method = "GET",
        data = null,
        params = null,
        headers = {},
        withCredentials = true
    } = options;

    // 构建完整 URL
    let fullUrl = url;
    if (params) {
        const queryString = Object.entries(params)
            .map(([key, value]) => `${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
            .join("&");
        fullUrl += (url.includes("?") ? "&" : "?") + queryString;
    }

    // 构建请求头
    const requestHeaders = {
        "Content-Type": "application/json",
        ...headers
    };

    // 构建请求配置
    const requestConfig = {
        method,
        headers: requestHeaders,
        credentials: withCredentials ? "include" : "same-origin"
    };

    // 添加请求体
    if (data && ["POST", "PUT", "PATCH"].includes(method.toUpperCase())) {
        requestConfig.body = typeof data === "string" ? data : JSON.stringify(data);
    }

    try {
        const response = await fetch(fullUrl, requestConfig);

        if (!response.ok) {
            throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        return await response.json();
    } catch (error) {
        console.error("Request failed:", error);
        throw error;
    }
}

// 快捷方法
const get = (url, options = {}) => request({ url, method: "GET", ...options });
const post = (url, data, options = {}) => request({ url, method: "POST", data, ...options });
const put = (url, data, options = {}) => request({ url, method: "PUT", data, ...options });
const del = (url, options = {}) => request({ url, method: "DELETE", ...options });

// 消息提示 (简化版，需要 UI 框架提供)
const message = {
    error: (msg) => console.error(msg),
    success: (msg) => console.success(msg),
    info: (msg) => console.info(msg),
    warning: (msg) => console.warn(msg)
};

export { request as r, request, get, post, put, delete: del, message };
export default request;
