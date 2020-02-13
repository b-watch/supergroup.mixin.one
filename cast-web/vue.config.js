module.exports = {
  publicPath: process.env.NODE_ENV === 'production'
    ? '/a0/cast/'
    : '/',
  "transpileDependencies": [
    "vuetify"
  ]
}
