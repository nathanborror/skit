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

var handleColor = function(hash) {
  var hue_resolution = 359;
  var hue = parseInt(hash, 16) % hue_resolution;
  return 'hsl(' + hue + ', '
                          + (.7 * 100) + '%, '
                          + (.7 * 100)+'%)';
};

// HACK
$(function() {
  $('footer a').each(function() {
    var user = $(this);
    var color = handleColor(user.data('hash'));
    user.css('background', color);
  });
});
