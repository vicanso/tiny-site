<template lang="pug">
nav.nav
  ul(
    v-for="item in categories"
  )
    li: a(
      :class="$route.name == routeListFile && active == item ? 'active': ''"
      href="javascript:;"
      @click="go(item)"
    ) {{item}}
</template>

<style lang="sass" scoped>
@import '@/styles/const.sass'
.nav
  position: fixed    
  left: 0
  top: $MAIN_HEADER_HEIGHT
  bottom: 0
  width: $MAIN_NAV_WIDTH
  background-color: $COLOR_BLACK
  z-index: 9
ul
  $height: 50px
  margin: 0
  padding: 0
  list-style: none
  line-height: $height
  a
    display: block
    color: $COLOR_WHITE
    text-decoration: none
    padding-left: 15px
    border-left: 3px solid $COLOR_BLACK
    &:hover, &.active
      color: $COLOR_BLUE
      background-color: $COLOR_WHITE
      border-left-color: $COLOR_BLUE
</style>

<script>
import { mapState } from "vuex";
import { routeListFile } from "@/routes";
export default {
  name: "main-nav",
  data() {
    return {
      active: this.$route.params.category || "",
      routeListFile
    };
  },
  computed: {
    ...mapState({
      categories: ({ file }) => file.categories
    })
  },
  methods: {
    go(category) {
      this.active = category;
      this.$router.push({
        name: routeListFile,
        params: {
          category
        }
      });
    }
  }
};
</script>
