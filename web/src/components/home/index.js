import React from "react";
import { Card } from "antd";

import "./home.sass";

class Home extends React.Component {
  render() {
    return (
      <div className="Home">
        <Card title="图片设置">
          <p>
            默认生成的预览地址并没有针对性的设置，实际使用中可根据应用场景调整图片参数以适应各不同场景，可调整的参数如下：
          </p>
          <ul>
            <li>图片质量</li>
            <li>图片高度</li>
            <li>图片宽度</li>
            <li>图片类型</li>
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
              指定图片质量，并调整图片宽高以及质量：
              <span>01DM5TZXH920856WDJ1JHTZTJV-70-100-300.webp</span>
            </li>
          </ul>
        </Card>
      </div>
    );
  }
}

export default Home;
