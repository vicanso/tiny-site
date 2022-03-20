import { defineComponent } from "vue";
import ExLoginRegister from "../components/ExLoginRegister";

export default defineComponent({
  name: "RegisterView",
  render() {
    return <ExLoginRegister type="register" />;
  },
});
