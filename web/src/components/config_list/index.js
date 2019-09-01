import React from "react";
import { Tabs } from "antd";

import ConfigTable from "../config_table";
import "./config_list.sass";
const { TabPane } = Tabs;

class ConfigList extends React.Component {
  renderConfigs(params = {}) {
    return (
      <ConfigTable
        params={params}
        formatData={data => {
          if (!data) {
            return "";
          }
          if (data[0] === "{" || data[0] === "[") {
            return <pre>{JSON.stringify(JSON.parse(data), null, 2)}</pre>;
          }
          return data;
        }}
      />
    );
  }
  render() {
    return (
      <div className="ConfigList">
        <Tabs defaultActiveKey="all" animated={false}>
          <TabPane tab="所有配置" key="all">
            {this.renderConfigs()}
          </TabPane>
          <TabPane tab="当前有效配置" key="available">
            {this.renderConfigs({
              available: true
            })}
          </TabPane>
          <TabPane tab="当前失效配置" key="unavailable">
            {this.renderConfigs({
              unavailable: true
            })}
          </TabPane>
        </Tabs>
      </div>
    );
  }
}

export default ConfigList;
