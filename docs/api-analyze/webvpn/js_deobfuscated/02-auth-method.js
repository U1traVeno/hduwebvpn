/**
 * 文件: authMethod-CtTsQ3_a.js → 02-auth-method.js
 * 功能: 本地账号绑定相关 API
 */

import { r as request } from "./request.js";

function startBindLocalAccount(data) {
    return request({ url: "/api/access/auth-method-bind-local/start", method: "post", data: data });
}

function bindLocalAccount(data) {
    return request({ url: "/api/access/auth-method/bind-local-account", method: "post", data: data });
}

function startHalfAccountBind(data) {
    return request({ url: "/api/access/half-account-bind/start", method: "post", data: data });
}

function finishHalfAccountBind(data) {
    return request({ url: "/api/access/half-account-bind/finish", method: "post", data: data });
}

export {
    finishHalfAccountBind as finishHalfBind,
    bindLocalAccount as bindLocal,
    startBindLocalAccount,
    startHalfAccountBind as startHalfBind
};
