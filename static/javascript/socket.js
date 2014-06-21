if (window['WebSocket']) {
  window.SOCKET = new WebSocket("ws://localhost:8080/ws");
  window.SOCKET.onopen = function(e) {
    console.log('Connection opened.');
  };
  window.SOCKET.onclose = function(e) {
    console.log('Connection closed.');
  };
  window.SOCKET.onmessage = function(e) {
    console.log('onmessage');
  };
}

window.State = {
  history: [],

  goto: function(url, user) {
    this.history.push({'url': url});
    this.refresh();
  },

  back: function() {
    this.history.pop();
    this.refresh();
  },

  refresh: function() {
    var payload = JSON.stringify(this.current());
    window.SOCKET.send(payload);
    window.history.pushState({}, "", this.current().url);
  },

  current: function() {
    if (this.history.length == 0) {
      return {'url': '/'};
    }
    return this.history[this.history.length-1];
  }
};
