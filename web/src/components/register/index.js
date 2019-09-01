import LoginRegister from "../login_register";

class Register extends LoginRegister {
  constructor(props) {
    super(props);
    this.state.mode = this.registerMode;
  }
}

export default Register;
