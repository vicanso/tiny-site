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
        prop="width"
        width="80"
        label="宽度"
      )
      el-table-column(
        prop="height"
        width="80"
        label="高度"
      )
      el-table-column(
        prop="creator"
        width="120"
        label="上传者"
      )
      el-table-column(
        fixed="right"
        label="操作"
        width="220"
      )
        template(
          slot-scope="scope"
        )
          el-button(
            @click="preview(scope.row)"
            type="text"
            size="small"
          ) 图片预览
          el-button(
            @click="clip(scope.row)"
            type="text"
            size="small"
          ) 图片剪辑
          el-button(
            @click="copyUrl(scope.row)"
            type="text"
            size="small"
          ) 复制链接
    .pagination: el-pagination(
      layout="sizes, prev, pager, next"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
      :page-sizes="[10, 20, 30, 50]"
      :pageSize="limit" 
      :total="fileCount"
      :current-page="currentPage"
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
import { IMAGES_PREVIEW, IMAGES_CLIP } from "@/urls";
import { urlPrefix } from "@/config";
import { saveListFilePageSize, getListFilePageSize } from "@/helpers/storage";
import { copy } from "@/helpers/util";
export default {
  name: "list-file",
  data() {
    let limit = getListFilePageSize();
    if (!limit) {
      limit = 10;
    }
    return {
      skip: 0,
      limit,
      loading: false,
      tableHeight: 0,
      category: this.$route.params.category,
      order: "-createdAt",
      fields: "file,category,maxAge,createdAt,type,size,creator,width,height",
      currentFiles: null,
      currentPage: -1
    };
  },
  computed: {
    ...mapState({
      files: ({ file }) => file.list,
      fileCount: ({ file }) => file.count,
      fileURLPrefix: ({ file }) => file.urlPrefix
    })
  },
  methods: {
    ...mapActions(["fileList", "fileCacheRemove"]),
    reset() {
      this.fileCacheRemove();
      this.skip = 0;
      this.currentPage = 1;
      this.currentFiles = null;
    },
    preview(data) {
      const { file, type, width, height } = data;
      const h = this.$createElement;
      const url = urlPrefix + IMAGES_PREVIEW.replace(":file", file);
      const maxWidth = 300;
      let style = "display:block;margin:auto;";
      if (width && height) {
        let newWidth = width;
        let newHeight = height;
        if (newWidth > maxWidth) {
          newHeight = (maxWidth / newWidth) * newHeight;
          newWidth = maxWidth;
        }
        style += `width:${newWidth}px;`;
        style += `height:${newHeight}px;`;
      }
      this.$msgbox({
        title: "图片预览",
        message: h("div", null, [
          h(
            "img",
            {
              attrs: {
                src: `${url}-90-0-0.${type}`
              },
              style
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
      const prefix = this.fileURLPrefix || location.origin;
      copy(`${prefix}${url}-90-0-0.${type}`);
      this.$message("已复制图片链接");
    },
    clip(data) {
      const { file, type } = data;
      const url =
        urlPrefix +
        IMAGES_CLIP.replace(":file", file).replace(":clip", "center");
      const width = Number.parseInt(data.width / 2);
      const height = Number.parseInt(data.height / 2);
      const prefix = this.fileURLPrefix || location.origin;
      copy(`${prefix}${url}-90-${width}-${height}.${type}`);
      this.$message("已复制图片剪辑链接（截取居中部分）");
    },
    handleSizeChange(val) {
      saveListFilePageSize(val);
      this.limit = val;
      this.reset();
      this.fetch();
    },
    handleCurrentChange(page) {
      this.skip = this.limit * (page - 1);
      this.currentPage = page;
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
        this.currentFiles = this.files.slice(skip, skip + limit);
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
