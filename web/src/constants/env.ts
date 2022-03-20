const env = import.meta.env;

// isDevelopment 是否开发环境
export function isDevelopment(): boolean {
  return env.DEV;
}

// isProduction 是否生产环境
export function isProduction(): boolean {
  return env.PROD;
}
