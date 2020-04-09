import React from "react";
import {
  Spin,
  message,
  Icon,
  Card,
  Input,
  Row,
  Col,
  Select,
  Upload,
  Button
} from "antd";
import bytes from "bytes";
import debounce from "debounce";

import * as fileService from "../../services/file";
import * as imageService from "../../services/image";
import { copy } from "../../helpers/util";
import "./preview.sass";
import { FIELS_UPLOAD } from "../../urls";

const { Option } = Select;

const { Dragger } = Upload;

class Preview extends React.Component {
  state = {
    inited: false,
    id: 0,
    clientX: 0,
    optiming: false,
    previewImage: null,
    optimImageInfo: null,
    optimWidth: 0,
    optimHeight: 0,
    optimQuality: 0,
    optimType: ""
  };
  previewImageRef = React.createRef();
  constructor(props) {
    super(props);
    this.state.id = Number.parseInt(props.match.params.id);
    this.debounceUpdateOptimParams = debounce((...args) => {
      this.updateOptimParams(...args);
    }, 1000);
    this._onPaste = e => {
      if (e.clipboardData && e.clipboardData.getData) {
        const url = e.clipboardData.getData("Text");
        if (url.indexOf("http://") === 0 || url.indexOf("https://") === 0) {
          // 通过URL获取图片预览
          this.downloadOriginalImage(url);
          return;
        }
        message.error(`仅支持http(s)形式，"${url}"不符合`);
      }
    };
    window.document.addEventListener("paste", this._onPaste);
  }
  componentWillUnmount() {
    window.document.removeEventListener("paste", this._onPaste);
  }
  getFileName() {
    const {
      optimQuality,
      optimWidth,
      optimHeight,
      previewImage,
      optimType
    } = this.state;
    const file = `${
      previewImage.name
    }-${optimQuality}-${optimWidth}-${optimHeight}.${optimType ||
      previewImage.type}`;
    return file;
  }
  getParams() {
    const {
      optimQuality,
      optimWidth,
      optimHeight,
      previewImage,
      optimType
    } = this.state;
    return {
      base64: previewImage.data,
      type: optimType || previewImage.type,
      sourceType: previewImage.type,
      quality: optimQuality,
      width: optimWidth,
      height: optimHeight
    };
  }
  updateOptimParams(key, value) {
    const { previewImage, optimQuality } = this.state;
    if (Number.isInteger(value) && value < 0) {
      message.error("参数不能小于0");
      return;
    }

    const update = {
      optimImageInfo: null
    };
    // 如果是转换为webp图片
    if (key === "optimType" && value === "webp") {
      // 原图片为jpeg，而且未选择转换质量，则设置为80来转换
      if (previewImage.type === "jpeg" && optimQuality === 0) {
        update.optimQuality = 80;
      }
    }

    update[key] = value;
    this.setState(update, () => {
      this.optimImage();
    });
  }
  download() {
    const { optimImageInfo } = this.state;
    const name = this.getFileName();
    const url = `data:image/${optimImageInfo.type};base64,${optimImageInfo.data}`;
    let a = document.createElement("a");
    let event = new MouseEvent("click");
    a.download = name;
    a.href = url;
    a.dispatchEvent(event);
  }
  async fetchConfig() {
    try {
      const data = await imageService.getConfig();
      this.setState({
        imageConfig: data
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  async optimImage() {
    const { optiming, id } = this.state;
    if (optiming) {
      return;
    }
    const file = this.getFileName();
    this.setState({
      optiming: true
    });
    try {
      let params = file;
      let fn = imageService.optim;
      if (!id) {
        params = this.getParams();
        fn = imageService.optimFromData;
      }
      const data = await fn(params);
      this.setState({
        optimImageInfo: data
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        optiming: false
      });
    }
  }
  handleMouseMove(e) {
    const { previewImage } = this.state;
    // 如果非预览，则直接返回
    if (!previewImage) {
      return;
    }
    const clientX = e.clientX;
    if (clientX % 3 === 0) {
      this.setState({
        clientX
      });
    }
  }
  async componentDidMount() {
    const { id } = this.state;
    try {
      this.fetchConfig();
      // 如果id为0，则表示直接上传文件预览
      if (!id) {
        this.setState({
          inited: true
        });
        return;
      }
      const data = await fileService.getByID(id, {
        fields: "*"
      });
      this.setState(
        {
          previewImage: data,
          inited: true
        },
        () => {
          this.optimImage();
        }
      );
    } catch (err) {
      message.error(err.message);
    }
  }
  async downloadOriginalImage(fileURL) {
    try {
      const data = await fileService.downloadFile(fileURL);
      this.setOrginalImage(data, data.name);
    } catch (err) {
      message.error("下载文件失败，" + err.message);
    }
  }
  setOrginalImage(response, name) {
    const self = this;

    let type = response.type;
    if (type === "jpg") {
      type = "jpeg";
    }
    const img = new Image();
    img.onload = () => {
      self.setState(
        {
          previewImage: {
            name,
            size: response.size,
            data: response.data,
            type,
            width: img.width,
            height: img.height
          }
        },
        () => {
          self.optimImage();
        }
      );
    };
    img.src = `data:image/${type};base64,${response.data}`;
  }
  renderUpload() {
    const self = this;
    const props = {
      muttiple: false,
      action: FIELS_UPLOAD,
      onChange(info) {
        const { status, response, name } = info.file;
        if (status === "done") {
          self.setOrginalImage(response, name);
        } else if (status === "error") {
          message.error("上传文件失败");
        }
      }
    };
    return (
      <Dragger {...props}>
        <p className="ant-upload-drag-icon">
          <Icon type="inbox" />
        </p>
        <p className="ant-upload-text">
          点击或拖动文件至此区域上传，支持复制URL后在此页面粘贴即可上传
        </p>
        <p className="ant-upload-hint">支持PNG与JEPG图片</p>
      </Dragger>
    );
  }
  renderPreview() {
    const {
      id,
      optiming,
      previewImage,
      optimImageInfo,
      optimWidth,
      optimType,
      optimHeight,
      optimQuality,
      imageConfig,
      clientX
    } = this.state;
    if (!previewImage) {
      return this.renderUpload();
    }
    // 高度-(60 + 50) 60为顶部高度，50为底部工具栏高度
    // 宽度-25
    const { innerHeight, innerWidth } = window;
    const maxWidth = innerWidth - 25;
    const maxHeight = innerHeight - 110;
    let originalWidth = previewImage.width;
    let originalHeight = previewImage.height;
    let currentWidth = 0;
    let currentHeight = 0;
    // 宽高都有设置
    if (optimWidth && optimHeight) {
      currentWidth = optimWidth;
      currentHeight = optimHeight;
    } else if (optimHeight) {
      // 只设置高度，宽度自适应
      currentHeight = optimHeight;
      currentWidth = (originalWidth / originalHeight) * currentHeight;
    } else if (optimWidth) {
      // 只设置宽度，高度自适应
      currentWidth = optimWidth;
      currentHeight = (originalHeight / originalWidth) * currentWidth;
    } else {
      currentWidth = originalWidth;
      currentHeight = originalHeight;
    }
    const wTimes = currentWidth / maxWidth;
    const hTimes = currentHeight / maxHeight;
    // 如果宽度或高度超过显示区域
    if (wTimes > 1 || hTimes > 1) {
      if (wTimes > hTimes) {
        currentHeight *= maxWidth / currentWidth;
        currentWidth = maxWidth;
      } else {
        currentWidth *= maxHeight / currentHeight;
        currentHeight = maxHeight;
      }
    }

    const close = (
      <a
        href="/close"
        className="close"
        onClick={e => {
          e.preventDefault();
          this.props.history.goBack();
        }}
      >
        <Icon type="close-square" />
      </a>
    );
    const tips = [];
    if (!previewImage) {
      tips.push("正在加载原图片数据");
    }
    if (!optimImageInfo) {
      tips.push("正在加载优化图片数据");
    }
    if (tips.length !== 0) {
      tips.push("请稍候...");
    }
    if (previewImage && optimImageInfo) {
      const originalSize = previewImage.size;
      const optimSize = optimImageInfo.size;
      tips.push(
        `优化后，文件大小${bytes(originalSize)} => ${bytes(
          optimSize
        )}，减少${100 - Math.floor((100 * optimSize) / originalSize)}%数据量`
      );
    }
    const title = `图片预览：${previewImage.name}`;
    const marginLeft = (maxWidth - currentWidth) / 2;
    const imgStyle = {
      position: "relative",
      marginLeft: `${marginLeft}px`,
      width: `${currentWidth}px`,
      height: `${currentHeight}px`,
      backgroundSize: "100% 100%",
      backgroundImage: `url("data:image/${previewImage.type};base64,${previewImage.data}")`
    };
    const halfOffset = currentWidth / 2;
    let leftValue = halfOffset;
    if (clientX && this.previewImageRef.current) {
      leftValue = clientX - this.previewImageRef.current.offsetLeft;
    }
    if (leftValue <= 0) {
      leftValue = halfOffset;
    } else if (leftValue >= currentWidth) {
      leftValue = halfOffset;
    }
    const optimStyle = {
      borderLeft: "1px solid #fff",
      position: "absolute",
      top: "0px",
      right: "0px",
      bottom: "0px",
      backgroundPosition: "right center",
      backgroundSize: "auto 100%",
      // 还要送去左边框的1px
      left: `${leftValue - 1}px`
    };
    if (optimImageInfo) {
      optimStyle.backgroundImage = `url("data:image/${optimImageInfo.type};base64,${optimImageInfo.data}")`;
    } else {
      optimStyle.backgroundColor = "rgba(255, 255, 255, 0.4)";
    }
    const file = this.getFileName();
    const url = imageConfig.url
      .replace(":zone", previewImage.zone)
      .replace(":file", file);
    const colList = [];
    if (!optiming) {
      colList.push(
        <Col key="qualityCol" span={3}>
          <Input
            defaultValue={optimQuality}
            addonBefore="图片质量："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimQuality",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="widthCol" span={3}>
          <Input
            defaultValue={optimWidth}
            addonBefore="图片宽度："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimWidth",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="heightCol" span={3}>
          <Input
            defaultValue={optimHeight}
            addonBefore="图片高度："
            type="number"
            onChange={e => {
              this.debounceUpdateOptimParams(
                "optimHeight",
                e.target.valueAsNumber || 0
              );
            }}
          />
        </Col>
      );
      colList.push(
        <Col key="typeCol" span={2}>
          <Select
            placeholder="图片类型"
            style={{
              width: "100%"
            }}
            defaultValue={optimType || previewImage.type}
            onChange={value => {
              this.updateOptimParams("optimType", value);
            }}
          >
            <Option value="png">PNG</Option>
            <Option value="jpeg">JPEG</Option>
            <Option value="webp">WEBP</Option>
          </Select>
        </Col>
      );
      if (id) {
        colList.push(
          <Col key="urlCol" span={5}>
            <Input
              readOnly
              addonBefore="图片地址："
              addonAfter={
                <Icon
                  title="点击复制图片地址"
                  onClick={e => {
                    copy(url, e.target);
                    message.info("已成功复制图片地址");
                  }}
                  type="copy"
                />
              }
              defaultValue={url}
            />
          </Col>
        );
      } else {
        colList.push(
          <Col key="baseCol" span={2}>
            <Button
              onClick={(e) => {
                copy(`data:image/${previewImage.type};base64,${previewImage.data}`, e.target);
                message.info("已成功复制图片base64数据");
              }}
              style={{
                width: "100%"
              }}
            >
              base64 
            </Button>
          </Col>,
          <Col key="urlCol" span={3}>
          <Button
            onClick={() => {
              this.download();
            }}
            style={{
              width: "100%"
            }}
          >
            下载图片
          </Button>
        </Col>
        );
      }
    }

    return (
      <Card title={title} extra={close} size="small" className="previewWrapper">
        <div className="content">
          <div
            className="imgWrapper"
            style={imgStyle}
            ref={this.previewImageRef}
          >
            <div style={optimStyle}></div>
            <div className="imgOriginal">原图</div>
            <div className="imgOptim">预览图</div>
          </div>
          <Row className="functions" gutter={12}>
            {colList}
            <Col span={8} className="tips">
              <Icon
                type="info-circle"
                style={{
                  marginRight: "3px"
                }}
              />
              {tips.join("，")}
            </Col>
          </Row>
        </div>
      </Card>
    );
  }
  render() {
    const { inited } = this.state;
    return (
      <div className="Preview" onMouseMove={this.handleMouseMove.bind(this)}>
        {!inited && (
          <div className="loadingWrapper">
            <Spin />
          </div>
        )}
        {inited && this.renderPreview()}
      </div>
    );
  }
}

export default Preview;
