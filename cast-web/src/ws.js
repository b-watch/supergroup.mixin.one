import { Base64 } from "js-base64";

const ws = {
  connect: function(url, opts) {
    this.websocket = new WebSocket(url);
    this.websocket.onopen = opts.onopen || this._onopen;
    this.websocket.onclose = opts.onclose || this._onclose;
    this.websocket.onmessage = evt => {
      const proc = opts.onmessage || this._onmessage;
      const msg = JSON.parse(evt.data);
      proc(msg);
    };
    this.websocket.onerror = opts.onerror || this._onerror;
  },

  _onopen: function() {
    console.log("websocket rocks");
  },

  _onclose: function() {
    console.log("disconnect");
  },

  _onmessage: function(evt) {
    console.log(evt.data);
    this.websocket.close();
  },

  _onerror: function(evt) {
    console.error(evt.data);
  }
};
export default ws;
