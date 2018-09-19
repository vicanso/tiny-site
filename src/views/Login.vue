<template lang="pug">
.login
  h4 登 录
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
    el-form-item
      el-button.submit(
        type="primary"
        @click="submit"
      ) 登录
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
  margin-top: -200px
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
export default {
  name: "login",
  data() {
    return {
      form: {}
    };
  },
  methods: {
    ...mapActions(["userLogin"]),
    async submit() {
      const { account, password } = this.form;
      if (!account || !password) {
        this.xError("用户名与密码不能为空");
        return;
      }
      const close = this.xLoading();
      try {
        await this.userLogin({
          account,
          password
        });
        this.$router.back();
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    }
  }
};
</script>
