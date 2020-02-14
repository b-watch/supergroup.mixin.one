module.exports = {
  publicPath: process.env.NODE_ENV === 'production'
    ? '/$PUB_PATH$/'
    : '/',
  devServer: {
    disableHostCheck: true,
  },
  configureWebpack: {
    resolve: {
      alias: {
        'vue$': 'vue/dist/vue.esm.js'
      }
    }
  }
}
