<template lang="pug">
.home
  h5 the right one is the best
  p 随着接入终端的种类繁多，以往一张图片所有终端使用的暴力做法已跟不上潮流，我们更倾向于不同的终端选择更适合的图片（尺寸、质量甚至于更优的图片格式），Tiny主要用于对图片做压缩或尺寸调整以及格式转换，更简单的方式适配更多的场景。
  ul
    li 图片名称的组成为：ID-质量-宽度-高度.后缀，图片质量根据应用场景来选择，一般而就言90算是比较好的选择。对于webp，如果需要无损压缩，则需要设置为0
    li 如果终端支持webp格式，尽量选择webp格式，其压缩质量更高（可以搜索各平台如何判断是否支持webp）
    li 对于不同移动终端，根据分辨率选择最合适的尺寸，尽可能不大图小用
    li 使用CDN回源机制，提升性能，避免相同图片需要多次处理

</template>

<style lang="sass" scoped>
@import "@/styles/const.sass"
.home
  padding: 30px 80px
  color: $COLOR_BLACK
  ul
    font-size: 14px
h5
  font-family: Oleo Script,cursive
  font-size: 32px
  text-align: center
  text-transform: uppercase
</style>


<script>
export default {
  name: "home",
  components: {}
};
</script>
