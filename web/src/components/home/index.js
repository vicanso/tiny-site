import React from "react";
import { Card } from "antd";
import { Link } from "react-router-dom";

import { PREVIEW_PATH } from "../../paths";

import "./home.sass";

class Home extends React.Component {
  render() {
    return (
      <div className="Home">
        <Card title="图片设置">
          <h3>游客体验</h3>
          <p>
            可直接上传图片，选择需要压缩的质量与格式，直接下载保存文件。
            <Link to={PREVIEW_PATH.replace(":id", 0)}>立即体验</Link>
          </p>

          <h3>测试账号/密码：tiny/123456</h3>
          <p>
            默认生成的预览地址并没有针对性的设置，实际使用中可根据应用场景调整图片参数以适应各不同场景，参数按顺序添加至文件名中，可仅指定前置的参数，如需要指定图片宽度，则必须指定图片质量、高度，后续的参数则可忽略，参数如下：
          </p>
          <ul>
            <li>图片质量：为0表示使用默认质量，JPEG(80)，PNG(90)</li>
            <li>图片高度：为0表示原图片高度</li>
            <li>图片宽度：为0表示原图片宽度</li>
            <li>图片类型：支持JPEG, PNG以及WEBP</li>
            <li>
              裁剪类型：按照指定的宽高对图片做裁剪，支持九种裁剪方式，1(left
              top) 2(top center) 3(right top) 4(left center) 5(center center)
              6(right center) 7(left bottom) 8(bottom center) 9(right
              bottom)，此参数在预览功能中并未提供选择，建议在终端中需动态裁剪显示区域时使用
            </li>
          </ul>
          <p>
            上述参数通过文件名来指定，格式如下：
            <span>名称-质量-宽度-高度.类型</span>
            ，质量、宽度以及高度为可选参数，需要注意，如果要指定后面的参数，则前面的参数也必须指定，可使用默认值0，如jpeg图片
            <span>01DM5TZXH920856WDJ1JHTZTJV</span>，以下为几种场景场景：
          </p>
          <ul>
            <li>
              转换图片为webp：<span>01DM5TZXH920856WDJ1JHTZTJV.webp</span>
            </li>
            <li>
              等比例调整图片宽度为80px：
              <span>01DM5TZXH920856WDJ1JHTZTJV-0-80-0.jpeg</span>
            </li>
            <li>
              压缩图片质量为30：<span>01DM5TZXH920856WDJ1JHTZTJV-30.jpeg</span>
            </li>
            <li>
              指定图片质量，并调整图片宽高以及图片类型：
              <span>01DM5TZXH920856WDJ1JHTZTJV-70-100-300.webp</span>
            </li>
            <li>
              指定图片质量，并做居中剪切以及图片类型：
              <span>01DM5TZXH920856WDJ1JHTZTJV-70-100-300-5.webp</span>
            </li>
          </ul>
        </Card>
      </div>
    );
  }
}

export default Home;
