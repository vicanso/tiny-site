import React from "react";
import {
  Upload,
  Icon,
  message,
  Form,
  Input,
  Row,
  Col,
  Select,
  Button,
  Card,
  Spin
} from "antd";

import "./file_handler.sass";
import * as fileService from "../../services/file";
import { getQueryParams } from "../../helpers/util";
import { FIELS_UPLOAD } from "../../urls";

const { Dragger } = Upload;
const { Option } = Select;
const { TextArea } = Input;
const addMode = "add";
const updateMode = "update";

class FileHandler extends React.Component {
  state = {
    mode: "",
    loading: false,
    zone: 0,
    zoneName: "",
    height: 0,
    width: 0,
    size: 0,
    name: "",
    type: "",
    data: "",
    maxAgeUnit: "h",
    maxAge: 24,
    description: ""
  };
  constructor(props) {
    super(props);
    const { state } = this;
    state.zone = Number.parseInt(props.match.params.fileZoneID);
    const { search } = props.location;
    state.zoneName = getQueryParams(search, "name");
    const fileID = getQueryParams(search, "fileID");
    if (fileID) {
      state.fileID = Number.parseInt(fileID);
      state.mode = updateMode;
    } else {
      state.mode = addMode;
    }
  }
  async componentDidMount() {
    const { fileID } = this.state;
    if (!fileID) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      const data = await fileService.getByID(fileID, {
        fields: "*"
      });
      const arr = data.maxAge.split(/(\d+)/);
      this.setState({
        height: data.height,
        width: data.width,
        size: data.size,
        name: data.name,
        type: data.type,
        data: data.data,
        maxAge: Number.parseInt(arr[1]),
        maxAgeUnit: arr[2],
        description: data.description
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  reset() {
    this.setState({
      name: "",
      data: "",
      description: ""
    });
  }
  async handleSubmit(e) {
    e.preventDefault();
    const {
      mode,
      fileID,
      zone,
      name,
      description,
      maxAge,
      maxAgeUnit,
      type,
      width,
      height,
      data
    } = this.state;
    if (!name || !maxAge || !description) {
      message.error("名称、缓存有效期以及描述均不能为空");
      return;
    }
    const params = {
      zone,
      name,
      description,
      maxAge: `${maxAge}${maxAgeUnit}`,
      type,
      width,
      height,
      data
    };
    this.setState({
      loading: true
    });
    try {
      if (mode === addMode) {
        await fileService.saveFile(params);
        message.info("新建文件成功");
        this.reset();
      } else {
        delete params.zone;
        delete params.name;
        params.id = fileID;
        await fileService.updateFile(params);
        message.info("更新文件成功");
        this.props.history.goBack();
      }
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderUpload() {
    const { data, type } = this.state;
    const previewStyle = {};
    if (data) {
      previewStyle[
        "backgroundImage"
      ] = `url(data:image/${type};base64,${data})`;
    }
    return (
      <Dragger
        name="file"
        action={FIELS_UPLOAD}
        showUploadList={false}
        onChange={info => {
          const { status, response, name } = info.file;
          if (status === "done") {
            message.info(`${name}上传成功`);
            const data = {
              height: response.height,
              width: response.width,
              size: response.size,
              type: response.type,
              data: response.data
            };

            if (!this.state.name) {
              data.name = response.name;
            }
            this.setState(data);
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
  renderForm() {
    const { name, description, data, maxAge, mode } = this.state;
    // 未上传文件时不需要展示
    if (!data) {
      return;
    }
    const maxAgeUnits = (
      <Select
        defaultValue="h"
        style={{
          width: 60
        }}
        onChange={value => {
          this.setState({
            maxAgeUnit: value
          });
        }}
      >
        <Option value="h">时</Option>
        <Option value="m">分</Option>
        <Option value="s">秒</Option>
      </Select>
    );
    return (
      <Form onSubmit={this.handleSubmit.bind(this)} className="form">
        <Row gutter={12}>
          <Col span={12}>
            <Form.Item label="名称">
              <Input
                type="text"
                placeholder="请输入图片名称"
                defaultValue={name}
                disabled={mode === updateMode}
                onChange={e => {
                  this.setState({
                    name: e.target.value
                  });
                }}
              />
            </Form.Item>
          </Col>
          <Col span={12}>
            <Form.Item label="缓存有效期">
              <Input
                type="number"
                placeholder="请输入该图片的缓存有效期"
                defaultValue={maxAge}
                onChange={e => {
                  this.setState({
                    maxAge: e.target.valueAsNumber
                  });
                }}
                addonAfter={maxAgeUnits}
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Form.Item label="描述">
              <TextArea
                defaultValue={description}
                autosize={{ minRows: 3, maxRows: 6 }}
                type="textarea"
                placeholder="请输入该图片的描述"
                onChange={e => {
                  this.setState({
                    description: e.target.value
                  });
                }}
              />
            </Form.Item>
          </Col>
          <Col span={24}>
            <Button className="submit" type="primary" htmlType="submit">
              {mode === addMode ? "新建" : "更新"}
            </Button>
          </Col>
        </Row>
      </Form>
    );
  }
  render() {
    const { loading, zoneName } = this.state;
    const title = `请填写相关信息(${zoneName || ""})`;
    return (
      <div className="FileHandler">
        <Spin spinning={loading}>
          <Card size="small" title={title} className="form">
            {this.renderUpload()}
            {this.renderForm()}
          </Card>
        </Spin>
      </div>
    );
  }
}

export default FileHandler;
