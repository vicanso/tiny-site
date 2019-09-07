import React from "react";
import { Upload, Icon, message } from "antd";

import "./file_handler.sass";

const { Dragger } = Upload;

class FileHandler extends React.Component {
  state = {
    height: 0,
    width: 0,
    size: 0,
    name: "",
    type: "",
    data: ""
  };
  renderUpload() {
    const { data } = this.state;
    const previewStyle = {};
    if (data) {
      previewStyle["backgroundImage"] = `url(${data})`;
    }
    return (
      <Dragger
        name="file"
        action="/files/v1/upload"
        showUploadList={false}
        onChange={info => {
          const { status, response, name, type } = info.file;
          if (status === "done") {
            message.info(`${name}上传成功`);
            this.setState({
              height: response.height,
              width: response.width,
              size: response.size,
              name: response.name,
              type: response.type,
              data: `data:${type};base64,${response.data}`
            });
          } else if (status === "error") {
            message.error(`${name}上传失败，${response.message}`);
          }
        }}
      >
        <div className="preview" style={previewStyle}>
          <p className="ant-upload-drag-icon">
            <Icon type="inbox" />
          </p>
          <p className="ant-upload-text">点击或拖动文件至此区域上传</p>
          <p className="ant-upload-hint">
            每次仅支持上传一个文件，再次上传会覆盖之前的文件，图片上传成功后，此背景图会展示为预览图。
          </p>
        </div>
      </Dragger>
    );
  }
  render() {
    return <div className="FileHandler">{this.renderUpload()}</div>;
  }
}

export default FileHandler;
