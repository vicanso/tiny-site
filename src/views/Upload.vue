<template lang="pug">
.wrapper: .upload
  h3 上传图片
  el-upload.dragWrapper.tac(
    ref="upload"
    drag
    :action="upload"
    :on-error="handleError"
    :on-success="handleSuccess"
    :on-remove="handleRemove"
    :on-exceed="handleExceed"
    :limit="1"
  )
    i.el-icon-upload
    .el-upload__text
      | 拖动文件至此或
      em 点击上传
    .el-upload__tip(
      slot="tip"
    ) 图片必须小于1MB
  el-button.save(
    type="primary"
    @click="save"
  ) 保存
  el-form.form(
    ref="form"
    :model="form"
    label-width="80px"
  )
    el-form-item(
      label="分类"
    )
      el-select.select(
        v-model="form.category"
        placeholder="请选择图片分类"
      )
        el-option(
          v-for="item in categories"
          :key="item"
          :value="item"
        )
    el-form-item(
      label="图片类型"
    )
      el-select.select(
        v-model="form.fileType"
        placeholder="请选择图片类型"
      )
        el-option(
          v-for="item in fileTypeList"
          :key="item"
          :value="item"
        )
    el-form-item(
      label="缓存时长"
    )
      el-select.select(
        v-model="form.maxAge"
        placeholder="请选择缓存时长"
      )
        el-option(
          v-for="item in maxAgeList"
          :key="item.value"
          :label="item.name"
          :value="item.value"
        )
    el-form-item(
      label="自定义"
    )
      el-input.input(
        v-model="form.customCategory"
        placeholder="自定义新的分类"
      )
    .desc 
      p
        i.el-icon-info.mright5
        span 缓存时长对应HTTP响应头中的Cache-Control: max-age
      p 自定义分类优先于选择的分类

</template>
<style lang="sass" scoped>
@import "@/styles/const.sass"
$dragWidth: 400px
$padding: 15px
$uploadWidth: 800px
.wrapper
  position: fixed
  left: $MAIN_NAV_WIDTH
  top: $MAIN_HEADER_HEIGHT
  right: 0
  bottom: 0
.upload
  padding: 0 $padding 
  position: absolute
  top: 50%
  left: 50%
  margin-top: -250px
  width: $uploadWidth
  margin-left: -$uploadWidth / 2 
  background-color: $COLOR_WHITE
  border: $GRAY_BORDER
  border-radius: 2px
h3
  margin: 0  
  line-height: 4em
  text-indent: 1em
.dragWrapper
  width: $dragWidth
  position: absolute
  top: 2 * $padding
  right: 0
  z-index: 1
.select, .input
  width: 280px
.save
  position: absolute
  bottom: 15px
  right: $padding
  width: $dragWidth - 2 * $padding
  z-index: 1
.desc
  margin-bottom: 15px
p
  font-size: 13px
  color: $COLOR_DARK_GRAY
  padding: 3px 0
  margin: 0
  margin-left: 16px
  &:first-child
    margin-left: 0
</style>


<script>
import { mapState, mapActions } from "vuex";
import { FILES_UPLOAD, IMAGES_PREVIEW } from "@/urls";
import { urlPrefix } from "@/config";
export default {
  name: "upload",
  data() {
    return {
      upload: urlPrefix + FILES_UPLOAD,
      form: {
        id: "",
        fileType: ""
      },
      fileTypeList: ["jpeg", "png"],
      maxAgeList: [
        {
          name: "1分钟",
          value: "1m"
        },
        {
          name: "5分钟",
          value: "5m"
        },
        {
          name: "30分钟",
          value: "30m"
        },
        {
          name: "1小时",
          value: "1h"
        },
        {
          name: "3小时",
          value: "3h"
        },
        {
          name: "10小时",
          value: "10h"
        },
        {
          name: "1天",
          value: "24h"
        },
        {
          name: "7天",
          value: "168h"
        },
        {
          name: "1月",
          value: "720h"
        },
        {
          name: "1年",
          value: "8760h"
        }
      ]
    };
  },
  computed: {
    ...mapState({
      categories: ({ file }) => file.categories
    })
  },
  methods: {
    ...mapActions(["fileSave"]),
    reset() {
      this.form = {};
      this.$refs.upload.clearFiles();
    },
    handleSuccess(res) {
      this.form.id = res.id;
      let fileType = res.fileType;
      if (fileType == "jpg") {
        fileType = "jpeg";
      }
      this.form.fileType = fileType;
    },
    handleExceed() {
      this.xError("请先移除当前已上传文件");
    },
    handleError(err) {
      this.xError(err);
    },
    handleRemove() {
      this.reset();
    },
    async save() {
      const { fileType, id, customCategory, maxAge } = this.form;
      let category = this.form.category;
      if (customCategory) {
        category = customCategory;
      }
      if (!fileType || !id || !category || !maxAge) {
        this.xError("图片、图片类型、缓存时长与图片分类不能为空");
        return;
      }
      const close = this.xLoading();
      try {
        const res = await this.fileSave({
          id,
          fileType,
          category,
          maxAge
        });
        this.reset();
        const { file, size, type } = res.data;
        this.$alert(
          `已成功上传图片，该图片id为：${file}，大小：${size}字节`,
          "成功上传",
          {
            confirmButtonText: "预览",
            callback: action => {
              if (action === "cancel") {
                return;
              }
              const url = urlPrefix + IMAGES_PREVIEW.replace(":file", file);
              window.open(`${url}-90-0-0.${type}`);
            }
          }
        );
      } catch (err) {
        this.xError(err);
      } finally {
        close();
      }
    }
  }
};
</script>
