export default function byteSize (n) {
  const unit = ['B', 'KB', 'MB', 'GB']
  let v = n
  let i = 0
  while (v > 1024) {
    v /= 1024
    i++
  }
  return v.toFixed(2) + ' ' + unit[i]
}