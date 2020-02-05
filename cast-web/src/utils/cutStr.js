export default function (str, len) {
  if (str.length > len * 2) {
    const s1 = str.substr(0, 8)
    const s2 = str.substr(str.length-8, str.length)
    return s1 + '...' + s2
  }
  return str
}