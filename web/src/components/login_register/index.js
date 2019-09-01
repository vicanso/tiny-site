import React from "react";
import { Form, Input, Icon, Card, Button, message, Spin, Row, Col } from "antd";

import { sha256 } from "../../helpers/crypto";
import { generatePassword } from "../../helpers/util";
import "./login_register.sass";
import * as userService from "../../services/user";
import * as commonService from "../../services/common";

class LoginRegister extends React.Component {
  loginMode = "login";
  registerMode = "register";
  state = {
    submitting: false,
    account: "",
    password: "",
    token: "",
    captcha: "",
    mode: "",
    captchaID: "",
    captchaData: ""
  };
  componentDidMount() {
    this.getCaptcha();
  }
  async getCaptcha() {
    this.setState({
      captchaData: ""
    });
    try {
      const data = await commonService.getCaptcha();
      this.setState({
        captchaID: data.id,
        captchaData: `data:image/${data.type};base64,${data.data}`
      });
    } catch (err) {
      message.error(err.message);
    }
  }
  async handleSubmit(e) {
    const { history } = this.props;
    const { setUserInfo } = this.props;
    e.preventDefault();
    const { account, password, mode, token, captchaID, captcha } = this.state;

    if (!account || !password) {
      message.error("用户名与密码不能为空");
      return;
    }
    if (!captcha || !captchaID) {
      message.error("图形验证码不能为空");
      return;
    }
    const postData = {
      captcha: `${captchaID}:${captcha}`,
      account
    };
    let fn = userService.login;
    if (mode === this.loginMode) {
      if (!token) {
        message.error("Token为不能空");
        return;
      }
      postData.password = sha256(generatePassword(password) + token);
    } else {
      fn = userService.register;
      postData.password = generatePassword(password);
    }
    this.setState({
      submitting: true
    });
    try {
      const data = await fn(postData);
      if (setUserInfo) {
        setUserInfo({
          account: data.account || "",
          roles: data.roles
        });
      }
      if (history) {
        history.goBack();
      }
    } catch (err) {
      message.error(err.message);
      // 因为图形验证码只可以使用一次，因此失败自动刷新
      this.getCaptcha();
    } finally {
      this.setState({
        submitting: false
      });
    }
  }
  render() {
    const { mode, submitting, captchaData } = this.state;
    const title = mode === this.loginMode ? "登录" : "注册";
    return (
      <div className="LoginRegister">
        <Spin spinning={submitting}>
          <Card title={title}>
            <Form onSubmit={this.handleSubmit.bind(this)}>
              <Form.Item>
                <Input
                  autoFocus
                  prefix={
                    <Icon type="user" style={{ color: "rgba(0,0,0,.25)" }} />
                  }
                  onChange={e => {
                    this.setState({
                      account: e.target.value.trim()
                    });
                  }}
                  placeholder="用户名"
                />
              </Form.Item>
              <Form.Item>
                <Input
                  prefix={
                    <Icon type="lock" style={{ color: "rgba(0,0,0,.25)" }} />
                  }
                  type="password"
                  onChange={e => {
                    this.setState({
                      password: e.target.value.trim()
                    });
                  }}
                  autoComplete="off"
                  placeholder="密码"
                />
              </Form.Item>
              <Form.Item>
                <Row gutter={8}>
                  <Col span={20}>
                    <Input
                      placeholder="请输入图形验证码"
                      onChange={e => {
                        const v = e.target.value.trim();
                        if (v.length > 4) {
                          message.warn("请输入4位长度验证码");
                        }
                        this.setState({
                          captcha: v
                        });
                      }}
                    />
                  </Col>
                  <Col span={4}>
                    <a
                      className="captcha"
                      href="/"
                      onClick={e => {
                        e.preventDefault();
                        this.getCaptcha();
                      }}
                    >
                      {captchaData && (
                        <img height="38" src={captchaData} alt="captcha" />
                      )}
                    </a>
                  </Col>
                </Row>
              </Form.Item>
              <Button type="primary" htmlType="submit" className="submit">
                {title}
              </Button>
            </Form>
          </Card>
        </Spin>
      </div>
    );
  }
}

export default LoginRegister;
