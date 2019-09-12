# tiny-site

在各终端分辨率以及支持图片格式各有差异的现实下，开发者一般都会采用折衷方案（偷懒）：选用最高的分辨率，最通用的图片格式。对于小分辨率终端，过高的分辨率，显示时做缩小展示，浪费了带宽。对于支持更优图片的终端，没有在质量与流量取得更优的平衡。

此图片管理系统，依赖于[tiny](https://github.com/vicanso/tiny)对图片做优化处理，可生成`webp`, `png`与`jpeg`。简便的形式支持自定义图片质量与尺寸，搭配CDN可根据应用需要生成各类不同的图片。


## 使用步骤

### 初始化数据库

数据库使用`postgres`，可以使用docker启动相关的镜像并设置初始化数据库，其中账号密码可根据需要设置。

```bash
docker run \
  -p 5432:5432 \
  -e POSTGRES_USER=test \
  -e POSTGRES_PASSWORD=123456 \
  --restart=always \
  --name=postgres \
  -d postgres:alpine
```


