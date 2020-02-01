export default function (str) {
  const reg = /(http:\/\/|https:\/\/)((\w|=|\?|\.|\/|&|-)+)/g
  return str.replace(reg, function(url) {
    return `<a target="_blank" href=${url}>${url}</a>`
  });
}