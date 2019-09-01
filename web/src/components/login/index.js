import { message } from "antd";
import PropTypes from "prop-types";

import LoginRegister from "../login_register";
import * as userService from "../../services/user";

class Login extends LoginRegister {
  constructor(props) {
    super(props);
    this.state.mode = this.loginMode;
  }
  async componentWillMount() {
    try {
      const data = await userService.getLoginToken();
      this.setState({
        token: data.token
      });
    } catch (err) {
      message.error(err.message);
    }
  }
}

Login.propTypes = {
  setUserInfo: PropTypes.func.isRequired
};

export default Login;
