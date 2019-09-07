import React from "react";
import { Spin, message } from "antd";

import "./file_list.sass";
import * as fileService from "../../services/file";

class FileList extends React.Component {
  state = {
    loading: false,
    zone: 0,
    offset: 0,
    limit: 10,
    fields:
      "id,updatedAt,name,maxAge,zone,type,size,width,height,description,creator"
  };
  constructor(props) {
    super(props);
    this.state.zone = Number.parseInt(props.match.params.fileZoneID);
  }
  componentDidMount() {
    this.fetchFiles();
  }
  async fetchFiles() {
    const { loading, offset, limit, zone, fields } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });

    try {
      await fileService.list({
        offset,
        limit,
        zone,
        fields
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderList() {}
  render() {
    const { loading } = this.state;
    return (
      <div className="FileList">
        <Spin spinning={loading}>{this.renderList()}</Spin>
      </div>
    );
  }
}

export default FileList;
