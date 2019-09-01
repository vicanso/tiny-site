import React from "react";
import { Card, Spin, Input, message, Table, DatePicker } from "antd";
import moment from "moment";

import "./user_login_record_list.sass";
import { TIME_FORMAT } from "../../vars";
import { setBeginOfDay, setEndOfDay } from "../../helpers/util";
import * as userService from "../../services/user";

const { Search } = Input;
const { RangePicker } = DatePicker;

class UserLoginRecordList extends React.Component {
  state = {
    loginRecords: null,
    account: "",
    pagination: {
      current: 1,
      pageSize: 10,
      total: 0
    },
    begin: moment(),
    end: moment(),
    total: 0,
    loading: false
  };
  async search() {
    const { loading, account, pagination, begin, end } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      const offset = (pagination.current - 1) * pagination.pageSize;
      const data = await userService.listLoginRecords({
        begin: setBeginOfDay(begin).toISOString(),
        end: setEndOfDay(end).toISOString(),
        account,
        limit: pagination.pageSize,
        offset
      });
      const updateData = {
        loginRecords: data.logins
      };
      if (data.count >= 0) {
        updateData.pagination = Object.assign(
          { ...pagination },
          {
            total: data.count
          }
        );
      }
      this.setState(updateData);
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  renderTable() {
    const { loginRecords, pagination } = this.state;
    const columns = [
      {
        title: "用户名",
        dataIndex: "account",
        key: "account"
      },
      {
        title: "IP",
        dataIndex: "ip",
        key: "ip"
      },
      {
        title: "Session ID",
        dataIndex: "sessionId",
        key: "sessionId"
      },
      {
        title: "Track ID",
        dataIndex: "trackId",
        key: "trackId"
      },
      {
        title: "X-Forwarded-For",
        dataIndex: "xForwardedFor",
        key: "xForwardedFor"
      },
      {
        title: "定位",
        key: "location",
        render: (text, record) => {
          const arr = [];
          if (record.country) {
            arr.push(record.country);
          }
          if (record.province) {
            arr.push(record.province);
          }
          if (record.city) {
            arr.push(record.city);
          }
          return arr.join(" ");
        }
      },
      {
        title: "ISP",
        dataIndex: "isp",
        key: "isp"
      },
      {
        title: "登录时间",
        dataIndex: "createdAt",
        key: "createdAt",
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      }
    ];
    return (
      <Table
        rowKey={"id"}
        className="loginRecords"
        dataSource={loginRecords}
        columns={columns}
        pagination={pagination}
        onChange={pagination => {
          this.setState(
            {
              pagination: { ...pagination }
            },
            () => {
              this.search();
            }
          );
        }}
      />
    );
  }
  renderLoginRecordList() {
    const { loading, begin, end } = this.state;
    return (
      <div>
        <Card title="登录搜索" size="small">
          <Spin spinning={loading}>
            <div className="filter">
              <Search
                className="account"
                placeholder="请输入账号"
                onSearch={account => {
                  this.setState({
                    account
                  });
                  this.search();
                }}
                enterButton
              />
              <RangePicker
                className="dateRange"
                defaultValue={[begin, end]}
                onChange={dates => {
                  this.setState({
                    begin: dates[0],
                    end: dates[1]
                  });
                }}
              />
            </div>
          </Spin>
        </Card>
        {this.renderTable()}
      </div>
    );
  }
  render() {
    return (
      <div className="UserLoginRecordList">{this.renderLoginRecordList()}</div>
    );
  }
}

export default UserLoginRecordList;
