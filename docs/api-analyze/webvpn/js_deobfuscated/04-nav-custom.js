/**
 * 文件: navCustom-CgDbtGc9.js → 04-nav-custom.js
 * 功能: 导航自定义相关 API 和状态管理
 *
 * 注意: 导入 index-Cey6Kqla.js (Vue 框架) 和 request.js
 */

import { V as createStore } from "./index-Cey6Kqla.js";
import { r as request } from "./request.js";

function getSiteList() {
    return request({ url: "/api/access/nav/site-list", method: "get" });
}

function getFavoriteSites() {
    return request({ url: "/api/access/nav/favorite-sites", method: "post" });
}

function addToFavorites(data) {
    return request({ url: "/api/access/nav/add-to-favorites", method: "post", data: data });
}

function removeFromFavorites(data) {
    return request({ url: "/api/access/nav/remove-from-favorites", method: "post", data: data });
}

function getNavConfig() {
    return request({ url: "/api/access/nav/config", method: "get" });
}

function getImufeCode() {
    return request({ url: "/api/access/imufe/code", method: "post" });
}

const navCustomStore = createStore("navCustom", {
    state: () => ({
        customInfo: {
            banner: "",
            copyright: "",
            title: "网站导航",
            favicon: "",
            css: "",
            logo: "",
            siteTypeTagColor: "",
            userNameColor: "",
            userNamePriority: "username",
            personalCenterBgColor: "",
            data: { problemDescription: "" }
        }
    }),
    actions: {
        async getCustomInfo() {
            const { code, data: response } = await request({ url: "/api/access/nav/custom", method: "get" });
            if (code === 0 && response) {
                this.customInfo = { ...response };
            }
        }
    },
    getters: {
        getBanner: (state) => state.customInfo.banner,
        getCopyright: (state) => state.customInfo.copyright,
        getTitle: (state) => state.customInfo.title,
        getFavicon: (state) => state.customInfo.favicon,
        getCss: (state) => state.customInfo.css,
        getLogo: (state) => state.customInfo.logo,
        getSiteTypeTagColor: (state) => state.customInfo.siteTypeTagColor,
        getUserNameColor: (state) => state.customInfo.userNameColor,
        getUserNamePriority: (state) => state.customInfo.userNamePriority ? state.customInfo.userNamePriority : "username",
        getPersonalCenterBgColor: (state) => state.customInfo.personalCenterBgColor,
        getProblemDescription: (state) => {
            return state.customInfo.data?.problemDescription || "";
        }
    }
});

export {
    addToFavorites as addFavorite,
    getSiteList as siteList,
    getFavoriteSites as favoriteSites,
    getNavConfig as navConfig,
    getImufeCode as imufeCode,
    removeFromFavorites as removeFavorite,
    navCustomStore as navStore
};
