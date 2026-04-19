/**
 * 文件: authentication-DJLiEriA.js → 01-authentication.js
 * 功能: 获取认证方式列表
 */

import { r as request } from "./request.js";

function getAuthenticationList() {
    return request({ url: "/api/access/authentication/list", method: "get" });
}

function getAllAuthentication() {
    return request({ url: "/api/access/authentication/all", method: "get" });
}

export { getAuthenticationList as getAuthList, getAllAuthentication as getAllAuthMethods };
