
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
  item.data({'hash': data.hash, 'parent': data.parent, 'root': data.root, 'user': data.user});
  return item;
}

// Saves a new item and inserts it into the DOM

function handleSave(e) {
  e.preventDefault();
  var form = $(this);

  $.post('/i/save', form.serialize(), function(data) {
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
  var url = '/i/'+item.attr('id');
  var data = {
    'hash': item.data('hash'),
    'parent': item.data('parent'),
    'root': item.data('root'),
    'user': item.data('user'),
    'text': item.data('text')
  };

  var viewItem = menu.find('.ui-item-view');
  viewItem.attr('href', url);
  viewItem.data(data);

  var deleteItem = menu.find('.ui-item-delete');
  deleteItem.attr('href', url+'/delete');
  deleteItem.data(data);

  var editItem = menu.find('.ui-item-edit');
  editItem.attr('href', url+'/edit');
  editItem.data(data);

  var colorItem = menu.find('.ui-item-colors');
  colorItem.data(data);

  menu.show();

  if ((e.pageY + menu.height()) > window.innerHeight) {
    menu.css({'top': e.pageY - menu.height(), 'left': e.pageX});
  } else {
    menu.css({'top': e.pageY, 'left': e.pageX});
  }
}

function handleColor(e) {
  e.preventDefault();

  var color = $(e.target).data('color');
  var hash = $(this).data('hash');
  var parent = $(this).data('parent');
  var root = $(this).data('root');
  var user = $(this).data('user');
  var text = $(this).data('text');
  var target = $('#'+hash).find('> a');
  target.css('background-color', 'rgba('+color+',1)');

  // var data = 'text='+text+'&hash='+hash+'&parent='+parent+'&root='+root+'&user='+user+'&color='+color+'';


  var serialized = $.param({
    'text': text,
    'hash': hash,
    'parent': parent,
    'root': root,
    'user': user,
    'color': color
  }, true);

  $.post('/i/save', serialized, function(data) {
    console.log(data);
  });
}

// HACK
$(function() {
  var body = $('body');

  body.on('click', '.ui-item a', handleItemClick);
  body.on('contextmenu', '.ui-item a', handleContextMenu);
  body.on('submit', '.ui-item-form', handleSave);
  body.on('click', '.ui-item-delete', handleDelete);
  body.on('click', '.ui-item-edit', handleEdit);
  body.on('click', '.ui-item-colors', handleColor);

  $(document).on('click', function() {
    $('#menu').hide();
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
