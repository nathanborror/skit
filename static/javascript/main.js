/*
 * WebSocket support
 */
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
