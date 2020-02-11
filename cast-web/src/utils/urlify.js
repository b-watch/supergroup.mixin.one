export default function (str) {
  const reg =  /((http|ftp|https):\/\/)?[\w\-_]+(\.[\w\-_]+)+([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?/g
  return str.replace(reg, function(url) {
    return `<a target="_blank" href=${url}>${url}</a>`
  });
}