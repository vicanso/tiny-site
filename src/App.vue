<template lang="pug">
#app
  MainHeader
  MainNav
  transition
    .mainWrapper.clearfix
      router-view
</template>
<style lang="sass" src="@/styles/app.sass"></style>

<script>
import { mapActions, mapState } from "vuex";
import MainHeader from "@/components/MainHeader.vue";
import MainNav from "@/components/MainNav.vue";

export default {
  name: "app",
  components: {
    MainHeader,
    MainNav
  },
  data() {
    return {
      account: ""
    };
  },
  methods: {
    ...mapActions(["userGetInfo", "fileListCategory"])
  },
  computed: {
    ...mapState({
      userInfo: ({ user }) => user.info
    })
  },
  watch: {
    userInfo(cur) {
      let account = "";
      if (cur) {
        account = cur.account;
      }
      if (account != this.account) {
        this.account = account;
        if (account) {
          this.fileListCategory();
        }
      }
    }
  },
  async beforeMount() {
    const close = this.xLoading();
    try {
      await this.userGetInfo();
    } catch (err) {
      this.xError(err);
    } finally {
      close();
    }
  }
};
</script>
