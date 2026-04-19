/**
 * 文件: auth-Crs8kf_D.js → 06-auth-crs.js
 * 功能: 认证核心 API
 */

import { r as request } from "./request.js";

function startAuth(data) {
    return request({ url: "/api/access/auth/start", method: "post", data: data });
}

function finishAuth(data) {
    return request({ url: "/api/access/auth/finish", method: "post", data: data });
}

function authTFA(data) {
    return request({ url: "/api/access/auth/tfa", method: "post", data: data });
}

function getTFAConfig() {
    return request({ url: "/api/access/auth/tfa-config", method: "post" });
}

function getUserNoticeInfo() {
    return request({ url: "/api/access/auth/user-notice-info", method: "get" });
}

function consumeSession(params) {
    return request({ url: "/api/access/auth/consume-session", method: "get", params: params });
}

function resetPassword(data) {
    return request({ url: "/api/access/auth/reset-password", method: "post", data: data });
}

function wechatLog(data) {
    return request({ url: "/api/access/auth/wechat-log", method: "post", data: data });
}

function getSessionToken(data) {
    return request({ url: "/api/access/auth/session-token", method: "post", data: data });
}

export {
    authTFA as doTFA,
    getTFAConfig as tfaConfig,
    startAuth as beginAuth,
    finishAuth as completeAuth,
    consumeSession as sessionConsume,
    wechatLog as wechatLogApi,
    getUserNoticeInfo as noticeInfo,
    resetPassword as passwordReset,
    getSessionToken as sessionToken
};
