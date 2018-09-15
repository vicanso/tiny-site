<template lang="pug">
.listFile(
  ref="listFile"
)
  .tableWrapper(
    v-if="!loading"
  )
    el-table.table(
      :data="currentFiles"
      :height="tableHeight"
      border
      stripe
    )
      el-table-column(
        prop="createdAt"
        label="创建时间"
        width="160"
      )
      el-table-column(
        prop="file"
        label="ID"
      )
      el-table-column(
        prop="maxAge"
        width="90"
        label="缓存时长"
      )
      el-table-column(
        prop="size"
        width="120"
        label="大小"
      )
      el-table-column(
        prop="creator"
        label="上传者"
      )
      el-table-column(
        fixed="right"
        label="操作"
        width="140"
      )
        template(
          slot-scope="scope"
        )
          el-button(
            @click="preview(scope.row)"
            type="text"
            size="small"
          ) 预览
          el-button(
            @click="copyUrl(scope.row)"
            type="text"
            size="small"
          ) 复制链接
    .pagination: el-pagination(
      layout="sizes, prev, pager, next"
      @size-change="handleSizeChange"
      :page-sizes="[10, 20, 30, 50]"
      :pageSize="limit" 
      :total="fileCount"
    )
</template>

<style lang="sass">
@import "@/styles/const.sass"
.listFile
  position: fixed 
  top: $MAIN_HEADER_HEIGHT
  left: $MAIN_NAV_WIDTH
  right: 0
  bottom: 0
  padding: 15px
.table
  width: 100%
.pagination
  padding: 15px 0
  text-align: right
</style>


<script>
import { mapState, mapActions } from "vuex";
import { IMAGES_PREVIEW } from "@/urls";
import { urlPrefix } from "@/config";
export default {
  name: "list-file",
  data() {
    return {
      skip: 0,
      limit: 10,
      loading: false,
      tableHeight: 0,
      category: this.$route.params.category,
      order: "-createdAt",
      fields: "file,category,maxAge,createdAt,type,size,creator",
      currentFiles: null
    };
  },
  computed: {
    ...mapState({
      files: ({ file }) => file.list,
      fileCount: ({ file }) => file.count
    })
  },
  watch: {
    files() {
      const { skip, limit } = this;
      this.currentFiles = this.files.slice(skip, skip + limit);
    }
  },
  methods: {
    ...mapActions(["fileList"]),
    reset() {
      this.skip = 0;
      this.currentFiles = null;
    },
    preview(data) {
      const { file, type } = data;
      const h = this.$createElement;
      const url = urlPrefix + IMAGES_PREVIEW.replace(":file", file);
      this.$msgbox({
        title: "图片预览",
        message: h("div", null, [
          h(
            "img",
            {
              attrs: {
                src: `${url}-90-0-0.${type}`
              },
              style: "display:block;margin:auto"
            },
            ""
          )
        ])
      })
        .then(() => {})
        .catch(() => {});
    },
    copyUrl(data) {
      const { file, type } = data;
      if (!document.execCommand) {
        this.$message({
          message: "很抱歉，该浏览器不支持复制",
          type: "warning"
        });
        return;
      }
      const url = urlPrefix + IMAGES_PREVIEW.replace(":file", file);
      // 来源自：https://juejin.im/post/5a94f8eff265da4e9b593c29
      const input = document.createElement("input");
      input.setAttribute("readonly", "readonly");
      input.setAttribute("value", `${location.origin}${url}-90-0-0.${type}`);
      document.body.appendChild(input);
      input.select();
      input.setSelectionRange(0, 9999);
      document.execCommand("copy");
      this.$message("已复制图片链接");
      document.body.removeChild(input);
    },
    handleSizeChange(val) {
      this.limit = val;
      this.reset();
      this.fetch();
    },
    async fetch() {
      const { skip, limit, order, fields, category } = this;
      const close = this.xLoading();
      this.loading = true;
      try {
        await this.fileList({
          skip,
          limit,
          order,
          category,
          fields
        });
      } catch (err) {
        this.xError(err);
      } finally {
        this.loading = false;
        close();
      }
    }
  },
  beforeRouteUpdate(to, from, next) {
    this.reset();
    this.category = to.params.category;
    this.fetch();
    next();
  },
  beforeMount() {
    this.fetch();
  },
  mounted() {
    this.tableHeight = this.$refs.listFile.clientHeight - 80;
  }
};
</script>
