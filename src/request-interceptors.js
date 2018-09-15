import request from "axios";

import { urlPrefix, timeout } from "@/config";

request.interceptors.request.use(config => {
  if (!config.timeout) {
    config.timeout = timeout;
  }
  config.url = `${urlPrefix}${config.url}`;
  return config;
});
