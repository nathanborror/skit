
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

  $('.ux-focus').focus();
});
