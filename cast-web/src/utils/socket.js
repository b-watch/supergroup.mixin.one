function execute (cb, ...args) {
  cb && cb(...args)
}

export default class Socket { 
  constructor () {
    this.MAX_RETRY_NUM = 3
    this.INTERVAL = 1000
    this.times = 0
    this.callbacks = {}
    this.connecting = false
  }

  connect (url, opts) {
    if (this.connecting) { return }
    this.connecting = true
    this.callbacks = opts
    execute(this.callbacks.onconnect)
    this.websocket = new WebSocket(url);
    this.websocket.onopen = this._onopen.bind(this);
    this.websocket.onclose = this._onclose.bind(this);
    this.websocket.onmessage = this._onmessage.bind(this);
    this.websocket.onerror = this._onerror.bind(this);
    this.url = url
  }

  reconnect () {
    if (this.times < this.MAX_RETRY_NUM) {
      this.times += 1
      this.connect(this.url, this.callbacks)
    } else {
      this.times = 0
      execute(this.callbacks.onfail)
    }
  }

  _onopen () {
    execute(this.callbacks.onconnected);
  }

  _onmessage (evt) {
    const msg = JSON.parse(evt.data)
    execute(this.callbacks.onmessage, msg)
  }

  _onerror () {
    setTimeout(() => { this.reconnect() }, this.INTERVAL);
  }

  _onclose () {
    setTimeout(() => { this.reconnect()}, this.INTERVAL);
  }
}