# tiny-site

图片优化管理系统，依赖于[tiny](https://github.com/vicanso/tiny)对图片做优化处理，可生成`webp`, `png`与`jpeg`。支持自定义图片质量与尺寸，搭配CDN可根据应用需要生成各类不同的图片。

管理后面的前端部分代码在`site`分支中单独开发管理，生成编译后的文件放在`assets`目录中与程序一起打包应用。


## static

将`site`分支中将dist中编译后的代码复制至`assets`目录中（不需要.map），之后再执行`packr`命令将静态文件打包。

## postgres

`tiny-site`使用的数据库为`postgres`，初始化数据库的步骤如下：

- 使用docker启动postgres

```bash
docker run --name tiny-site \
  -p 5432:5432 \
  --restart=always \
  -v /opt/postgres/tiny-site:/var/lib/postgresql/data \
  -d postgres:alpine
```

- 创建数据库用户等

```bash
docker exec -it tiny-site sh

# 切换用户
su - postgres

psql

# 创建密码
\password

CREATE USER 账号 WITH PASSWORD '密码';
CREATE DATABASE tiny-site OWNER 账号;
GRANT ALL PRIVILEGES ON DATABASE tiny-site to 账号;
```

复制configs/default.yml中的db.uri配置至production中，修改数据库连接中的user/password字段

## production.yml

需要自定义的几个配置项：

```yaml
urlPrefix: /api
db:
  uri: postgres://user:pwd@127.0.0.1:5432/tiny-site?connect_timeout=5&sslmode=disable
tiny:
  address: 127.0.0.1:3016
session:
  keys:
    - aVOHyH
```

## 启动tiny-site

```bash
docker run -d \
  --restart=always \
  -e GO_ENV=production \
  -e CONFIG=/configs \
  -v /opt/tiny-site/production.yml:/configs/production.yml \
  -p 8080:8080 \
  --name tiny-site \
  vicanso/tiny-site
```

## test

```bash
GO_ENV=test VIPER_INIT_TEST=true go test -race -coverprofile=test.out ./... && go tool cover --html=test.out
```