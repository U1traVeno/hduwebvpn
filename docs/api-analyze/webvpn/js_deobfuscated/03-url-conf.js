/**
 * 文件: urlConf-D5f8gSXt.js → 03-url-conf.js
 * 功能: URL 配置文件
 */

const baseCallbackUrl = `${window.location.origin}/callback`;

const callbackUrls = {
    oauth2CallbackUrl: `${baseCallbackUrl}/oauth2`,
    casCallbackUrl: `${baseCallbackUrl}/cas`,
    larkCallbackUrl: `${baseCallbackUrl}/lark`,
    wechatEnCallbackUrl: `${baseCallbackUrl}/wechat-enterprise`,
    wechatEnPcCallbackUrl: `${baseCallbackUrl}/wechat-en-pc`,
    dingTalkCallbackUrl: `${baseCallbackUrl}/ding-talk`,
    wechatOfficialRzkcPcCallbackUrl: `${baseCallbackUrl}/wechat-official-rzkc-pc`,
    wechatOfficialRzkcCallbackUrl: `${baseCallbackUrl}/wechat-official-rzkc`,
    alipayPcCallbackUrl: `${baseCallbackUrl}/alipay-pc`,
    alipayCallbackUrl: `${baseCallbackUrl}/alipay`,
    ecardCallbackUrl: `${baseCallbackUrl}/ecard`
};

const getReturnUrl = async () => {
    const defaultReturnUrl = `${window.location.origin}/site-nav/`;
    return `${localStorage.getItem("returnUrl") ? localStorage.getItem("returnUrl") : defaultReturnUrl}${localStorage.getItem("returnUrlHash") ? localStorage.getItem("returnUrlHash") : ""}`;
};

export { callbackUrls as callbackConfig, getReturnUrl as getHomeUrl };
