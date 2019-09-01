import axios from "axios";

import { ROUTERS } from "../urls";

export async function list(params) {
  const { data } = await axios.get(ROUTERS, {
    params
  });
  return data;
}
