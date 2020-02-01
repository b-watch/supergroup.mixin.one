import hashcode from '@/utils/hashcode'
import colors from '@/assets/colors'

const nameColors = colors.names

export default function (str, map) {
  if (map[str]) { return map[str] }
  return nameColors[hashcode(str) % nameColors.length]
}