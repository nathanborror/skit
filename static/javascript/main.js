
function handleMessage(data) {
  var parent = $('#'+data.item.hash);
  replaceItems(data.items, parent);
}

// Kicks off an xhr request to pull down child items, if they exist.
// Or hides items if they're already showing.

function handleItemClick(e) {
  e.preventDefault();
  var item = $(this).parent();

  if (item.hasClass('ui-item-expanded')) {
    item.find('.ui-item').remove();
    item.find('.ui-item-form').remove();
    item.removeClass('ui-item-expanded');
    window.SOCKET.unsubscribe(this.pathname);
  } else {
    $.ajax({
      'url': this.href,
      'success': addItems
    });
    window.SOCKET.subscribe(this.pathname, handleMessage);
  }
}

// Appends child items to the clicked item.
// It also creates a new input field for that parent item.

function addItems(data) {
  var parent = $('#'+data.item.hash);
  var anchor = parent.find('> a');
  parent.addClass('ui-item-expanded');

  for (var i=0; i<data.items.length; i++) {
    var item = getItemHTML(data.items[i], 'ui-item-child');
    anchor.after(item);
  }

  // Copy form to items
  var form = $('#root-form').clone().appendTo(parent);
  form.attr('id', '');
  var parent_field = form.find('[name="parent"]');
  var root_field = form.find('[name="root"]');
  parent_field.val(data.item.hash)
  root_field.val(data.item.root)
  form.find('input[name="text"]').focus();
}

function replaceItems(items, parent) {
  parent.find('.ui-item').remove();
  var anchor = parent.find('> a');
  for (var i=0; i<items.length; i++) {
    var item = getItemHTML(items[i], 'ui-item-child');
    anchor.after(item);
  }
}

// Returns HTML necessary to render an item

function getItemHTML(data, extraClass) {
  var item = $('<div class="ui-item '+extraClass+'" id="'+data.hash+'"><a href="/i/'+data.hash+'">'+data.text+'</a></div>');
  item.data({'hash': data.hash, 'parent': data.parent, 'root': data.root});
  return item;
}

// Saves a new item and inserts it into the DOM

function handleSave(e) {
  e.preventDefault();
  var form = $(this);

  $.post('/i/save', $(this).serialize(), function(data) {
    if (form.parent().prop('tagName') == 'ARTICLE') {
      var parent = $('.ui-root-items')
      var item = getItemHTML(data.item);
      parent.append(item);
    } else {
      var parent = $('#'+data.item.parent);
      var item = getItemHTML(data.item, 'ui-item-child');
      form.before(item);
    }
    window.SOCKET.request('/i/'+data.item.parent);
  }.bind(this));

  form.find('input[name="text"]').val("");
}

function handleEdit(e) {
  e.preventDefault();
  alert("Not implemented yet :(");
}

function handleDelete(e) {
  e.preventDefault();

  var target = $(this);
  $.post('/i/'+target.data('hash')+'/delete', function(data) {
    var item = $('#'+target.data('hash'));
    item.remove();
    window.SOCKET.request('/i/'+target.data('parent'));
  });
}

// Shows or hides a contextual menu

function handleContextMenu(e) {
  e.preventDefault();
  var menu = $('#menu');
  var item = $(this).parent();

  menu.find('.ui-item-share').attr('href', '/i/'+item.attr('id'));
  menu.find('.ui-item-share').data({'hash': item.data('hash'), 'parent': item.data('parent'), 'root': item.data('root')});

  menu.find('.ui-item-delete').attr('href', '/i/'+item.attr('id')+'/delete');
  menu.find('.ui-item-delete').data({'hash': item.data('hash'), 'parent': item.data('parent'), 'root': item.data('root')});

  menu.find('.ui-item-edit').attr('href', '/i/'+item.attr('id')+'/edit');
  menu.find('.ui-item-edit').data({'hash': item.data('hash'), 'parent': item.data('parent'), 'root': item.data('root')});

  menu.show();
  menu.css({'top': e.pageY, 'left': e.pageX});
}

function handleColor(hash) {
  var hue_resolution = 359;
  var hue = parseInt(hash, 16) % hue_resolution;
  return 'hsl(' + hue + ', ' + (.7 * 100) + '%, ' + (.7 * 100)+'%)';
}

// HACK
$(function() {
  $('body').on('click', '.ui-item-list .ui-item a', handleItemClick);
  $('body').on('contextmenu', '.ui-item-list .ui-item a', handleContextMenu);
  $('body').on('submit', '.ui-item-list .ui-item-form', handleSave);
  $('body').on('click', '.ui-item-list .ui-item-delete', handleDelete);
  $('body').on('click', '.ui-item-list .ui-item-edit', handleEdit);

  $(document).on('click', function() {
    $('#menu').hide();
  });

  $('footer a').each(function() {
    var user = $(this);
    var color = handleColor(user.data('hash'));
    user.css('background', color);
  });

  $('.ux-focus').focus();

  window.SOCKET.onclose = function(e) {
    var alert = $("#alert");
    alert.append("<p>Oops! You've been disconnected. Refresh to fix this.</p>");
  }

  window.SOCKET.onopen = function(e) {
    var p = $("#alert p");
    p.remove();
  }
});
