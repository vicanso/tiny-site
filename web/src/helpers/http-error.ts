class HTTPError extends Error {
  // http状态码
  status: number;
  // 是否异常
  exception?: boolean;
  // 出错信息
  message: string;
  // 出错分类
  category?: string;
  // 出错代码，如果某个出错需要单独处理，可定义唯一的出错码
  code?: string;
  // 子错误
  errs?: HTTPError[];
  // 其它一些额外的信息
  extra?: Record<string, unknown>;
  constructor(status: number, message: string) {
    super(message);
    this.status = status;
    this.message = message;
  }
}

export default HTTPError;
