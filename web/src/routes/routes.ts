import { Component } from "vue";

export interface Router {
  path: string;
  name: string;
  component: Component | Promise<Component>;
}

export const names = {
  home: "home",
  profile: "profile",
  login: "logon",
  register: "register",
  users: "users",
  logins: "logins",
  userTrackers: "userTrackers",
  httpErrors: "httpErrors",
  requests: "requests",
  mockTime: "mockTime",
  configs: "configs",
  blockIPs: "blockIPs",
  signedKeys: "signedKeys",
  routerMocks: "routerMocks",
  routerConcurrencies: "routerConcurrencies",
  requestConcurrencies: "requestConcurrencies",
  caches: "caches",
  emails: "emails",
  httpServerInterceptors: "httpServerInterceptors",
};

export const routes: Router[] = [
  {
    path: "/",
    name: names.home,
    component: () => import("../views/Home"),
  },
  {
    path: "/profile",
    name: names.profile,
    component: () => import("../views/Profile"),
  },
  {
    path: "/login",
    name: names.login,
    component: () => import("../views/Login"),
  },
  {
    path: "/register",
    name: names.register,
    component: () => import("../views/Register"),
  },
  {
    path: "/users",
    name: names.users,
    component: () => import("../views/Users"),
  },
  {
    path: "/logins",
    name: names.logins,
    component: () => import("../views/stats/Logins"),
  },
  {
    path: "/user-trackers",
    name: names.userTrackers,
    component: () => import("../views/stats/UserTrackers"),
  },
  {
    path: "/http-errors",
    name: names.httpErrors,
    component: () => import("../views/stats/HTTPErrors"),
  },
  {
    path: "/requests",
    name: names.requests,
    component: () => import("../views/stats/Requests"),
  },
  {
    path: "/mock-time",
    name: names.mockTime,
    component: () => import("../views/configurations/MockTime"),
  },
  {
    path: "/configs",
    name: names.configs,
    component: () => import("../views/configurations/Configs"),
  },
  {
    path: "/block-ips",
    name: names.blockIPs,
    component: () => import("../views/configurations/BlockIPs"),
  },
  {
    path: "/signed-keys",
    name: names.signedKeys,
    component: () => import("../views/configurations/SignedKeys"),
  },
  {
    path: "/router-mocks",
    name: names.routerMocks,
    component: () => import("../views/configurations/RouterMocks"),
  },
  {
    path: "/router-concurrencies",
    name: names.routerConcurrencies,
    component: () => import("../views/configurations/RouterConcurrencies"),
  },
  {
    path: "/request-concurrencies",
    name: names.requestConcurrencies,
    component: () => import("../views/configurations/RequestConcurrencies"),
  },
  {
    path: "/caches",
    name: names.caches,
    component: () => import("../views/Caches"),
  },
  {
    path: "/emails",
    name: names.emails,
    component: () => import("../views/configurations/Emails"),
  },
  {
    path: "/http-server-interceptors",
    name: names.httpServerInterceptors,
    component: () => import("../views/configurations/HTTPServerInterceptors"),
  },
];
