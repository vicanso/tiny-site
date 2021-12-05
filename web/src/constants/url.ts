// 用户相关url
// 用户信息
export const USERS_ME = "/users/v1/me";
// 用户详细信息
export const USERS_ME_DETAIL = "/users/v1/detail";
// 用户登录
export const USERS_LOGIN = "/users/v1/me/login";
export const USERS_INNER_LOGIN = "/users/inner/v1/me/login";
// 用户列表
export const USERS = "/users/v1";
// 用户登录记录
export const USERS_LOGINS = "/users/v1/login-records";
export const USERS_ID = "/users/v1/:id";

// 通用接口相关url
// 图形验证码
export const COMMONS_CAPTCHA = "/commons/captcha";
// 路由列表
export const COMMONS_ROUTERS = "/commons/routers";
// HTTP性能指标统计
export const COMMONS_HTTP_STATS = "/commons/http-stats";

// flux相关查询
// 用户行为日志列表
export const FLUXES_TRACKERS = "/fluxes/v1/trackers";
// http出错列表
export const FLUXES_HTTP_ERRORS = "/fluxes/v1/http-errors";
// 客户端上传的action日志列表
export const FLUXES_ACTIONS = "/fluxes/v1/actions";
// 后端HTTP调用列表
export const FLUXES_REQUESTS = "/fluxes/v1/requests";
// tag value列表
export const FLUXES_TAG_VALUES = "/fluxes/v1/tag-values/:measurement/:tag";
// flux单条记录查询
export const FLUXES_FIND_ONE = "/fluxes/v1/one/:measurement";

// 系统配置相关url
// 配置列表
export const CONFIGS = "/configurations/v1";
// 根据ID查询或更新配置
export const CONFIGS_ID = "/configurations/v1/:id";
// 当前有效配置
export const CONFIGS_CURRENT_VALID = "/configurations/v1/current-valid";

// 管理员相关接口
export const ADMINS_CACHE_ID = "/@admin/v1/caches/:key";
