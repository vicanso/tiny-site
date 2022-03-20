import {
  NButton,
  NCard,
  NForm,
  NFormItem,
  NGrid,
  NGridItem,
  NInput,
  useMessage,
} from "naive-ui";
import { css } from "@linaria/core";
import { defineComponent, onMounted, ref, PropType } from "vue";
import { showError, showWarning } from "../helpers/util";
import { commonGetCaptcha, commonGetEmptyCaptcha } from "../states/common";
import { userLogin, userRegister } from "../states/user";
import { goToHome, goToLogin } from "../routes";

const loginType = "login";

const cardClass = css`
  max-width: 640px;
  margin: 120px auto;
`;

const captchaClass = css`
  text-align: center;
  height: 40px;
  cursor: pointer;
`;

const submitButtonClass = css`
  width: 100%;
`;
const captchaImageClass = css`
  height: 100%;
`;

export default defineComponent({
  name: "ExLoginRegister",
  props: {
    type: {
      type: String as PropType<"login" | "register">,
      default: loginType,
    },
  },
  setup(props) {
    const model = ref({
      account: "",
      password: "",
      captcha: "",
    });
    const message = useMessage();
    const isLogin = props.type === loginType;
    // 登录或注册中
    const processing = ref(false);
    const captchaData = ref(commonGetEmptyCaptcha());
    // 刷新图形验证码
    const refreshCaptcha = async () => {
      try {
        captchaData.value.id = "";
        const data = await commonGetCaptcha();
        captchaData.value = data;
      } catch (err) {
        showError(message, err);
      }
    };
    // 登录或注册
    const loginOrRegister = async () => {
      if (processing.value) {
        return;
      }
      const { account, password, captcha } = model.value;
      if (!account || !password || !captcha) {
        showWarning(message, "账号、密码及图形验证码均不能为空");
        return;
      }
      try {
        processing.value = true;
        const params = {
          account,
          password,
          captcha: `${captchaData.value.id}:${captcha}`,
        };
        if (isLogin) {
          await userLogin(params);
        } else {
          await userRegister(params);
        }
        // 跳转至首页
        if (isLogin) {
          goToHome(true);
        } else {
          goToLogin(true);
        }
      } catch (err) {
        showError(message, err);
        // 因为图形验证码仅一次有效
        // 如果出错，自动刷新图形验证码
        refreshCaptcha();
      } finally {
        processing.value = false;
      }
    };
    const handleEnterPress = (e: KeyboardEvent) => {
      if (e.code !== "Enter") {
        return;
      }
      loginOrRegister();
    };
    onMounted(() => {
      // 首次加载自动加载图形验证码
      refreshCaptcha();
    });

    return {
      refreshCaptcha,
      loginOrRegister,
      handleEnterPress,
      isLogin,
      captchaData,
      model,
      processing,
      rules: {
        account: {
          required: true,
          message: "账号不能为空",
          trigger: "blur",
        },
        password: {
          required: true,
          message: "密码不能为空",
          trigger: "blur",
        },
        captcha: {
          required: true,
          message: "图形验证码不能为空",
          trigger: "blur",
        },
      },
    };
  },
  render() {
    const size = "large";
    const { model, captchaData, rules, processing, isLogin, handleEnterPress } =
      this;
    const title = isLogin ? "用户登录" : "用户注册";
    const btnText = isLogin ? "登录" : "注册";
    return (
      <NCard title={title} class={cardClass}>
        <NForm
          labelWidth="100"
          labelAlign="right"
          labelPlacement="left"
          size={size}
          rules={rules}
          model={model}
        >
          <NFormItem label="账号：" path="account">
            <NInput
              placeholder="请输入账号"
              clearable
              autofocus
              onUpdateValue={(value) => {
                model.account = value;
              }}
            />
          </NFormItem>
          <NFormItem label="密码：" path="password">
            <NInput
              type="password"
              placeholder="请输入密码"
              clearable
              show-password-on="click"
              onUpdateValue={(value) => {
                model.password = value;
              }}
            />
          </NFormItem>
          <NFormItem label="验证码：" path="captcha">
            <NGrid cols={4}>
              <NGridItem span={3}>
                <NInput
                  placeholder="请输入验证码"
                  maxlength={4}
                  clearable
                  onKeyup={handleEnterPress}
                  onUpdateValue={(value) => {
                    model.captcha = value;
                  }}
                />
              </NGridItem>
              <NGridItem span={1}>
                <div class={captchaClass} onClick={this.refreshCaptcha}>
                  {captchaData.id && (
                    <img
                      class={captchaImageClass}
                      src={`data:image/${captchaData.type};base64,${captchaData.data}`}
                    />
                  )}
                  {!captchaData.id && <span>...</span>}
                </div>
              </NGridItem>
            </NGrid>
          </NFormItem>
          <NButton
            loading={processing}
            class={submitButtonClass}
            size={size}
            onClick={this.loginOrRegister}
          >
            {btnText}
          </NButton>
        </NForm>
      </NCard>
    );
  },
});
