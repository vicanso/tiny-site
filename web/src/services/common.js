import axios from "axios";

import { RANDOM_KEYS, CAPTCHA } from "../urls";

export async function getRandomKeys(params) {
  const { data } = await axios.get(RANDOM_KEYS, {
    params
  });
  return data;
}

export async function getCaptcha() {
  const { data } = await axios.get(CAPTCHA);
  return data;
}
