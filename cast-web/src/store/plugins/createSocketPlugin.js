import { WS_BASE_URL, SUP_MESSAGE_CAT } from '@/constants';

const MESSAGE_WS = WS_BASE_URL + "/messages";

export default function (socket) {
  return function (store) {
    const onmessage = function (msg) {
      if (SUP_MESSAGE_CAT.includes(msg.category)) {
        store.commit('message/ADD_MESSAGE')
      }
    }
    const onconnect = function () {
      store.commit('message/CONNECTING')
    }
    const onconnected = function () {
      store.commit('message/CONNECTED')
    }
    const onfail = function () {
      store.commit('message/CONNECT_FAILED')
    }
    socket.connect(MESSAGE_WS, { onmessage, onconnect, onconnected, onfail })
  }
}