import Vue from "vue";
import Router from "vue-router";
import Home from "@/views/Home.vue";
import Login from "@/views/Login.vue";
import Register from "@/views/Register.vue";
import Upload from "@/views/Upload.vue";
import ListFile from "@/views/ListFile.vue";
import {
  routeLogin,
  routeHome,
  routeRegister,
  routeUpload,
  routeListFile
} from "@/routes";

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: "/",
      name: routeHome,
      component: Home
    },
    {
      path: "/login",
      name: routeLogin,
      component: Login
    },
    {
      path: "/register",
      name: routeRegister,
      component: Register
    },
    {
      path: "/upload",
      name: routeUpload,
      component: Upload
    },
    {
      path: "/list-file/:category",
      name: routeListFile,
      component: ListFile
    }
  ]
});
