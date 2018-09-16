# tiny-site

图片优化管理系统，依赖于[tiny](https://github.com/vicanso/tiny)对图片做优化处理，可生成`webp`, `png`与`jpeg`。支持自定义图片质量与尺寸，搭配CDN可根据应用需要生成各类不同的图片。

管理后面的前端部分代码在`site`分支中单独开发管理，生成编译后的文件放在`assets`目录中与程序一起打包应用。

# docker

## docker build

## static

Create static asset's packr

```bash
packr
```

## test

```bash
GO_ENV=test VIPER_INIT_TEST=true go test -race -coverprofile=test.out ./... && go tool cover --html=test.out
```