<template lang="pug">
.home
  h5 the right one is the best
  h6 file-quality-width-height.ext
  p 随着接入终端的种类繁多，以往一张图片所有终端使用的暴力做法已跟不上潮流，我们更倾向于不同的终端选择更适合的图片（尺寸、质量甚至于更优的图片格式），Tiny主要用于对图片做压缩或尺寸调整以及格式转换，更简单的方式适配更多的场景。
  ul
    li 图片名称的组成为：ID-质量-宽度-高度.后缀，图片质量根据应用场景来选择，一般而就言90算是比较好的选择。对于webp，如果需要无损压缩，则需要设置为0
    li 如果终端支持webp格式，尽量选择webp格式，其压缩质量更高（可以搜索各平台如何判断是否支持webp）
    li 对于不同移动终端，根据分辨率选择最合适的尺寸，尽可能不大图小用
    li 使用CDN回源机制，提升性能，避免相同图片需要多次处理
    li 分类查询接口缓存1分钟，因此如果是新增加的分类需要等最长一分钟在列表中才会有
  
  p 注：我配置了七牛CDN回源，大家可以体验一下，如`http://tiny.aslant.site/api/images/v1/center/01CQPFKV5WYG8VD4QTVN5PYYJW-90-480-300.webp`，则使用`http://oidm8hv4x.qnssl.com/api/images/v1/center/01CQPFKV5WYG8VD4QTVN5PYYJW-90-480-300.webp`，在图片第一次使用时由于CDN无数据，需要从源头拉取会较慢。

</template>

<style lang="sass" scoped>
@import "@/styles/const.sass"
.home
  padding: 30px 80px
  color: $COLOR_BLACK
  ul, p
    font-size: 14px
h5, h6
  font-family: Oleo Script,cursive
  text-align: center
h5
  font-size: 32px
  text-transform: uppercase
h6
  font-size: 18px
</style>


<script>
export default {
  name: "home",
  components: {}
};
</script>
