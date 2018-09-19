<template lang="pug">
.login
  h4 注 册
  el-form.form(
    ref="form"
    :model="form"
    label-width="40px"
  )
    el-form-item(
      label="账号"
    )
      el-input(
        v-model="form.account"
        @keyup.enter.native="submit"
        autofocus
      )
    el-form-item(
      label="密码"
    ) 
      el-input(
        v-model="form.password"
        type="password"
        @keyup.enter.native="submit"
      )
    el-form-item(
      label="密码"
    ) 
      el-input(
        v-model="form.passwordConfirm"
        type="password"
        @keyup.enter.native="submit"
      )
    el-form-item
      el-button.submit(
        type="primary"
        @click="submit"
      ) 注册
</template>
<style lang="sass" scoped>
@import '@/styles/const.sass'

.login
  $width: 400px
  position: fixed
  width: $width
  left: 50%
  top: 50%
  margin-left: -$width / 2
  margin-top: -250px
  border: $GRAY_BORDER
  background-color: $COLOR_WHITE
h4
  margin: 0
  line-height: 2.5em
  background-color: $COLOR_BLACK
  color: $COLOR_WHITE
  padding: 5px 15px
.form
  padding: 20px
.submit
  width: 100%
</style>

<script>
import { mapActions } from "vuex";
import { routeLogin } from "@/routes";
export default {
  name: "register",
  data() {
    return {
      form: {}
    };
  },
  methods: {
    ...mapActions(["userRegister"]),
    async submit() {
      const { account, password, passwordConfirm } = this.form;
      if (!account || !password) {
        this.xError("用户名与密码不能为空");
        return;
      }
      if (password != passwordConfirm) {
        this.xError("两次输入的密码不相同");
        return;
      }
      const close = this.xLoading();
      try {
        await this.userRegister({
          account,
          password
        });
        this.$router.replace({
          name: routeLogin
        });
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    }
  }
};
</script>
