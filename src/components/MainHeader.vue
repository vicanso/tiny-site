<template lang="pug">
header.header
  .pullRight.mright10
    .functions(
      v-if="userInfo"
    )
      div(
        v-if="userInfo.account"
      )
        el-button.mleft10(
          type="text"
          @click="$router.push({name: routeUpload})"
        )
          i.el-icon-upload.mright5
          | 上传图片
        span.divide | 
        span.mright10
          i.el-icon-star-off.mright5
          | {{userInfo.account}}
        el-button(
          type="text"
          @click="userLogout()"
        )
          | 退出
          i.el-icon-back.mleft5
      div(
        v-else
      )
        el-button(
          type="text"
          @click="$router.push({name: routeLogin})"
        )
          i.el-icon-info.mright5
          | 登录
        el-button(
          type="text"
          @click="$router.push({name: routeRegister})"
        )
          i.el-icon-circle-plus.mright5
          | 注册
    i.el-icon-loading.font24(
      v-else
    )
  a.logo(
    href="javascript:;"
    @click="$router.push({name: routeHome})"
  ) Tiny
</template>

<style lang="sass" scoped>
@import '@/styles/const.sass'
.header
  position: fixed
  left: 0
  top: 0
  right: 0
  height: $MAIN_HEADER_HEIGHT
  background-color: $COLOR_WHITE
  border-bottom: $GRAY_BORDER
  z-index: 9
  line-height: $MAIN_HEADER_HEIGHT
  overflow: hidden
.logo
  width: $MAIN_NAV_WIDTH
  background: $COLOR_BLACK
  height: 100%
  color: $COLOR_WHITE
  text-indent: 2em
  display: block
  text-decoration: none
.divide
  margin: 0 20px
  color: $COLOR_DARK_GRAY
</style>

<script>
import { mapState, mapActions } from "vuex";
import { routeLogin, routeRegister, routeUpload, routeHome } from "@/routes";
export default {
  name: "main-header",
  data() {
    return {
      routeLogin,
      routeRegister,
      routeUpload,
      routeHome
    };
  },
  methods: {
    ...mapActions(["userLogout"])
  },
  computed: {
    ...mapState({
      userInfo: ({ user }) => user.info
    })
  }
};
</script>
