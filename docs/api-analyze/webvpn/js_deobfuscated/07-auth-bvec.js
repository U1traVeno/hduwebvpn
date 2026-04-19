/**
 * 文件: auth-Bvec_zE4.js → 07-auth-bvec.js
 * 功能: 认证 API 封装 (根据是否需要绑定本地账号选择不同流程)
 *
 * 导入已反混淆的文件:
 *   - 06-auth-crs.js (原 auth-Crs8kf_D.js)
 *   - 02-auth-method.js (原 authMethod-CtTsQ3_a.js)
 *   - user-V10Yxw5H.js (未反混淆，保持原名)
 */

import { d as completeAuth, c as beginAuth } from "./06-auth-crs.js";
import { b as bindLocal, n as startBindLocal } from "./02-auth-method.js";
import { u as loginStore } from "./user-V10Yxw5H.js";

const storeInstance = loginStore();

/**
 * 不需要绑定本地账号的认证流程
 */
const authWithoutLocal = async (credentials) => {
    return storeInstance.userInfo.needToBindLocalAccount
        ? await startBindLocal(credentials)
        : await beginAuth(credentials);
};

/**
 * 需要绑定本地账号的认证流程
 */
const authWithLocal = async (credentials) => {
    return storeInstance.userInfo.needToBindLocalAccount
        ? await bindLocal(credentials)
        : await completeAuth(credentials);
};

export {
    authWithoutLocal as authenticate,
    authWithLocal as authenticateWithBinding
};
