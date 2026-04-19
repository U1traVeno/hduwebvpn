/**
 * 文件: CasCallback-q97cqoLl.js → 08-cas-callback.js
 * 功能: CAS 回调处理组件 (Vue 组件)
 *
 * ⚠️ 这是基于源码逻辑的简化重构版本
 *
 * 导入已反混淆的文件:
 *   - 07-auth-bvec.js (原 auth-Bvec_zE4.js)
 *   - 05-user-api.js (原 user-cwdmpEuP.js)
 *   - 03-url-conf.js (原 urlConf-D5f8gSXt.js)
 *
 * 保持原名的文件 (未反混淆):
 *   - index-Cey6Kqla.js (Vue 框架)
 *   - user-V10Yxw5H.js (Pinia store)
 *   - util-DyruD4Ub.js (工具函数)
 *   - index-o4p_ZbnI.js (Vue 组件)
 */

import { a as authenticate } from "./07-auth-bvec.js";
import { _ as defineComponent, d as withDefaults, c as createElementVNode, o as createBaseVNode, f as Loading, h as ref, i as onMounted, B as message, k as useRoute, u as useRouter, j as toRefs } from "./index-Cey6Kqla.js";
import { u as loginStore } from "./user-V10Yxw5H.js";
import { g as generateDeviceId } from "./util-DyruD4Ub.js";
import { S as Spin } from "./index-o4p_ZbnI.js";
import { b as request } from "./request.js";

import "./06-auth-crs.js";
import "./02-auth-method.js";
import "./05-user-api.js";
import "./const-BHerojsC.js";
import "./03-url-conf.js";
import "./styleChecker-TFxV7Lzl.js";
import "./index-DXhXwGTY.js";
import "./initDefaultProps-_hAYIFu4.js";
import "./fade-BdgfkLw7.js";

const CasCallback = defineComponent({
    name: "CasCallback",

    setup() {
        const router = useRouter();
        const route = useRoute();
        const store = loginStore();
        const state = ref({ authLoading: false });

        onMounted(async () => {
            const userResult = await store.getUserInfo({ unbind: true });

            if (!(userResult?.data?.userId)) {
                if (state.value.authLoading = true, route.query.ticket) {
                    const ticket = decodeURIComponent(route.query.ticket);
                    const externalId = route.currentRoute.value.params.externalId;
                    const callbackUrl = window.location.href.split("?")[0];

                    const authPayload = {
                        externalId: externalId,
                        data: JSON.stringify({
                            callbackUrl: callbackUrl,
                            ticket: ticket,
                            deviceId: await generateDeviceId()
                        })
                    };

                    const doAuth = async (payload) => {
                        try {
                            const { code } = await authenticate(payload);
                            if (code === 0) {
                                state.value.authLoading = false;
                                store.getUserInfo();
                                console.log("userStore.showBindWechat", store.showBindWechat);
                            }
                        } finally {
                            state.value.authLoading = false;
                        }
                    };

                    await doAuth(authPayload);
                } else {
                    request.error("认证失败");
                }
            }
        });

        return {
            ...toRefs(state)
        };
    },

    render() {
        const SpinComponent = Spin;
        return createBaseVNode("div", { class: "callback-container" }, [
            createElementVNode(Loading, { tip: "授权中...", size: "large" })
        ]);
    }
});

export { CasCallback as default };
