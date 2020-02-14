module.exports = {
  publicPath: process.env.NODE_ENV === 'production'
    ? '/$PUB_PATH$/'
    : '/',
  "transpileDependencies": [
    "vuetify"
  ]
}
