import React from "react";
import {
  Row,
  Col,
  Input,
  Card,
  Select,
  Table,
  message,
  Spin,
  Button,
  Form
} from "antd";
import moment from "moment";

import "./user_list.sass";
import { TIME_FORMAT } from "../../vars";
import * as userService from "../../services/user";

const { Search } = Input;
const { Option } = Select;
const editMode = "edit";

const allRole = "all";
const roles = [allRole, "su", "admin"];

class UserList extends React.Component {
  state = {
    mode: "",
    keyword: "",
    role: "",
    current: null,
    loading: false,
    users: null,

    newRoles: null
  };
  async search() {
    const { loading, keyword, role } = this.state;
    if (loading) {
      return;
    }
    this.setState({
      loading: true
    });
    try {
      let filterRole = role;
      if (filterRole === allRole) {
        filterRole = "";
      }
      const data = await userService.list({
        limit: 20,
        keyword,
        role: filterRole
      });
      const { users } = data;
      users.forEach(item => {
        item.key = `${item.id}`;
      });
      this.setState({
        users
      });
    } catch (err) {
      message.error(err.message);
    } finally {
      this.setState({
        loading: false
      });
    }
  }
  async handleSubmit(e) {
    e.preventDefault();
    const { current, newRoles } = this.state;
    try {
      await userService.updateByID(current.id, {
        roles: newRoles
      });
      const users = this.state.users.slice(0);
      users.forEach(item => {
        if (item.id === current.id) {
          item.roles = newRoles;
        }
      });
      message.info("更新用户信息成功");
      this.setState({
        mode: "",
        users
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  renderTable() {
    const { users } = this.state;
    const columns = [
      {
        title: "用户名",
        dataIndex: "account",
        key: "account"
      },
      {
        title: "角色",
        dataIndex: "roles",
        key: "roles",
        render: roles => {
          if (roles) {
            return roles.join(",");
          }
          return;
        }
      },
      {
        title: "创建于",
        dataIndex: "createdAt",
        key: "createdAt",
        render: text => {
          if (!text) {
            return;
          }
          return moment(text).format(TIME_FORMAT);
        }
      },
      {
        title: "操作",
        key: "op",
        width: "100px",
        render: (text, record) => {
          return (
            <a
              href="/update"
              onClick={e => {
                e.preventDefault();
                this.setState({
                  newRoles: record.roles,
                  current: record,
                  mode: editMode
                });
              }}
            >
              更新
            </a>
          );
        }
      }
    ];
    return <Table className="users" dataSource={users} columns={columns} />;
  }
  renderRoles() {
    return roles.map(item => (
      <Option key={item} value={item}>
        {item}
      </Option>
    ));
  }
  renderUserList() {
    const { loading, mode } = this.state;
    if (mode === editMode) {
      return;
    }
    return (
      <div>
        <Card title="用户搜索" size="small">
          <Spin spinning={loading}>
            <div className="filter">
              <Select
                defaultValue={allRole}
                className="roles"
                placeholder="请选择用户角色"
                onChange={value => {
                  this.setState({
                    role: value
                  });
                }}
              >
                {this.renderRoles()}
              </Select>
              <Search
                className="keyword"
                placeholder="请输入关键字"
                onSearch={keyword => {
                  this.setState({
                    keyword
                  });
                  this.search();
                }}
                enterButton
              />
            </div>
          </Spin>
        </Card>
        {this.renderTable()}
      </div>
    );
  }
  renderEditor() {
    const { mode, current } = this.state;
    if (mode !== editMode) {
      return;
    }
    const colSpan = 12;
    return (
      <Card title="更新用户信息" size="small">
        <Form onSubmit={this.handleSubmit.bind(this)}>
          <Row gutter={24}>
            <Col span={colSpan}>
              <Form.Item label="用户名">
                <Input disabled defaultValue={current.account} />
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Form.Item label="用户角色">
                <Select
                  defaultValue={current.roles}
                  mode="multiple"
                  placeholder="请选择要添加的角色"
                  onChange={value => {
                    this.setState({
                      newRoles: value
                    });
                  }}
                >
                  {this.renderRoles()}
                </Select>
              </Form.Item>
            </Col>
            <Col span={colSpan}>
              <Button className="submit" type="primary" htmlType="submit">
                更新
              </Button>
            </Col>
            <Col span={colSpan}>
              <Button
                className="back"
                onClick={() => {
                  this.setState({
                    mode: ""
                  });
                }}
              >
                返回
              </Button>
            </Col>
          </Row>
        </Form>
      </Card>
    );
  }
  render() {
    return (
      <div className="UserList">
        {this.renderUserList()}
        {this.renderEditor()}
      </div>
    );
  }
}

export default UserList;
