const { defineConfig } = require('@vue/cli-service')
module.exports = defineConfig({
  transpileDependencies: true,
  
  // 根据生产环境设置正确的公共路径
  publicPath: process.env.NODE_ENV === 'production'
    ? '/web/'
    : '/',

  devServer: {
    proxy: {
      '/api': {
        target: process.env.VUE_APP_DEV_PROXY_TARGET, // 从环境变量读取代理目标地址
        changeOrigin: true,
        pathRewrite: {
          '^/api': '/v1/api'
        }
      }
    }
  }
})
