// 从请求的query中获取key对应的值
function getReqQuery(key) {
  return req.query[key];
}

// 修改请求头的query
function setReqQuery(key, value) {
  req.modifiedQuery = true;
  req.query[key] = value;
}

// 从请求数据中获取body中key对应的值
function getReqBody(key) {
  return req.body[key];
}

// 设置请求数据中的body的值
function setReqBody(key, value) {
  req.modifiedBody = true;
  req.body[key] = value;
}

// 设置响应数据
function setRespBody(key, value) {
  if (!resp.status) {
    resp.status = 200;
  }
  resp.body[key] = value;
}

// 设置响应状态码
function setRespStatus(status) {
  resp.status = status;
}

// 设置响应HTTP头
function setRespHeader(key, value) {
  resp.header[key] = value;
}