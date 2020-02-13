module.exports = {
  publicPath: process.env.NODE_ENV === 'production'
    ? '/a0/'
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
