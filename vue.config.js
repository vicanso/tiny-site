module.exports = {
  devServer: {
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true
      }
    }
  },
  baseUrl: process.env.NODE_ENV === 'production'
    ? '/static/'
    : '/'
};
