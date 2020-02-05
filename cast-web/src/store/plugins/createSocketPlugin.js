import { WS_BASE_URL } from '@/constants';

const MESSAGE_WS = WS_BASE_URL + "/messages";

export default function (socket) {
  return function (store) {
    const onmessage = function (msg) {
      store.commit('message/ADD_MESSAGE', msg)
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