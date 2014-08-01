
var handleMessage = function(data) {
  var item = $('#'+data.item.hash);
  if (item.length > 0) {
    ItemManager.updateChildren(data.items, item);
  } else {
    var root = $('.ui-root-items');
    ItemManager.updateItems(data.items, root);
  }
}

var ItemManager = {};

// AddItems Appends child items to the clicked item.
// It also creates a new input field for that parent item.
ItemManager.addItems = function(data) {
  var parent = $('#'+data.item.hash);
  var anchor = parent.find('> a');
  parent.addClass('ui-item-expanded');

  for (var i=0; i<data.items.length; i++) {
    var item = Item.html(data.items[i], 'ui-item-child');
    anchor.after(item);
  }

  // Copy form to items
  var form = $('#root-form').clone().appendTo(parent);
  form.attr('id', '');
  var parent_field = form.find('[name="parent"]');
  var root_field = form.find('[name="root"]');
  var color_field = form.find('[name="color"]');
  parent_field.val(data.item.hash)
  root_field.val(data.item.root)
  color_field.val(data.item.color)
  form.find('input[name="text"]').focus();
}

// UpdateItems checks all the items on screen and adds any missing from the
// current dataset.
ItemManager.updateItems = function(items, parent) {
  var root = $('.ui-root-items');
  var current = _.map(root.find('> .ui-item'), function(obj) {
    return obj.id;
  });

  var hashes = _.map(items, function(obj) {
    return obj.hash;
  });

  // Get the difference between the incoming hashes compared against
  // the existing hashes in the DOM.
  var diff = _.difference(hashes, current);

  // Insert any hashes that don't exist.
  for (var i=0; i<diff.length; i++) {
    var data = _.findWhere(items, {'hash': diff[i]});
    var item = Item.html(data);
    Item.insert(item, root);
  }

  // Check for hashes that exist in the DOM that shouldn't.
  var diff = _.difference(current, hashes);

  // Remove them from the DOM
  for (var i=0; i<diff.length; i++) {
    var item = $('#'+diff[i]);
    item.remove();
  }
}

// UpdateChildren checks all the child items on screen and adds any missing
// children from the current dataset.
ItemManager.updateChildren = function(items, parent) {
  parent.find('.ui-item').remove();
  var anchor = parent.find('> a');
  for (var i=0; i<items.length; i++) {
    var item = Item.html(items[i], 'ui-item-child');
    anchor.after(item);
  }
}

// Submit submits an item form.
ItemManager.submit = function(e) {
  e.preventDefault();
  var form = $(this);

  Item.save(form.serialize(), function(data) {
    if (form.parent().prop('tagName') == 'ARTICLE') {
      var parent = $('.ui-root-items')
      var item = Item.html(data.item);
      Item.insert(item, parent);
    } else {
      var item = Item.html(data.item, 'ui-item-child');
      form.before(item);
    }
  });

  form.find('input[name="text"]').val("");
}

// Color saves an assigned color to an item.
ItemManager.color = function(e) {
  e.preventDefault();

  var target = $(this);
  var item = $('#'+target.data('hash'));
  var color = $(e.target).data('color');

  item.find('> a').css('background-color', 'rgba('+color+',1)');

  target.data('color', color);
  var data = $.param(target.data(), true);
  Item.save(data);
}

var Item = {};

// HandleClick kicks off an xhr request to pull down child items, if they
// exist. Or hides items if they're already showing.
Item.handleClick = function(e) {
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
      'success': ItemManager.addItems
    });
    window.SOCKET.subscribe(this.pathname, handleMessage);
  }
}

// HTML returns HTML necessary to render an item.
Item.html = function(data, extraClass) {
  var item = $('<div class="ui-item '+extraClass+'" id="'+data.hash+'"><a href="/i/'+data.hash+'" style="background-color:rgba('+data.color+',.5); border-color:rgba('+data.color+',1);">'+data.text+'</a></div>');
  item.data({
    'hash': data.hash,
    'parent': data.parent,
    'root': data.root,
    'user': data.user,
    'text': data.text,
    'color': data.color
  });
  return item;
}

// Save saves a new item.
Item.save = function(data, complete) {
  $.post('/i/save', data, function(data) {
    if (complete) {
      complete(data);
    }

    window.SOCKET.request('/i/'+data.item.parent);
  }.bind(this));
}

// Edit edits an item.
Item.edit = function(e) {
  e.preventDefault();
  alert("Not implemented yet :(");
}

// Deletes removes an item.
Item.delete = function(e) {
  e.preventDefault();
  var target = $(this);
  $.post('/i/'+target.data('hash')+'/delete', function(data) {
    if (data.error) {
      console.log(data.error);
    } else {
      var item = $('#'+target.data('hash'));
      item.remove();
      window.SOCKET.request('/i/'+target.data('parent'));
    }
  });
}

// ContextMenu shows or hides a contextual menu for an item.
Item.contextMenu = function(e) {
  e.preventDefault();
  var menu = $('#menu');
  var item = $(this).parent();
  var url = '/i/'+item.attr('id');

  menu.find('.ui-item-view')
    .attr('href', url)
    .data(item.data());

  menu.find('.ui-item-delete')
    .attr('href', url+'/delete')
    .data(item.data());

  menu.find('.ui-item-edit')
    .attr('href', url+'/edit')
    .data(item.data());

  menu.find('.ui-item-colors')
    .data(item.data());

  menu.show();

  if ((e.pageY + menu.height()) > window.innerHeight) {
    menu.css({'top': e.pageY - menu.height(), 'left': e.pageX});
  } else {
    menu.css({'top': e.pageY, 'left': e.pageX});
  }
}

// Insert adds an item into a given parent node.
Item.insert = function(item, parent) {
  parent.prepend(item);
}

// HACK
$(function() {
  var body = $('body');

  body.on('click', '.ui-item a', Item.handleClick);
  body.on('contextmenu', '.ui-item a', Item.contextMenu);
  body.on('submit', '.ui-item-form', ItemManager.submit);
  body.on('click', '.ui-item-delete', Item.delete);
  body.on('click', '.ui-item-edit', Item.edit);
  body.on('click', '.ui-item-colors', ItemManager.color);

  $(document).on('click', function() {
    $('#menu').hide();
  });

  $('.ux-focus').focus();

  window.SOCKET.onclose = function(e) {
    var alert = $("#alert");
    alert.append("<p>Oops! You've been disconnected. <a href='javascript:location.reload();'>Reload</a> to fix this.</p>");
  }
  window.SOCKET.onopen = function(e) {
    var p = $("#alert p");
    p.remove();
    window.SOCKET.subscribe(window.location.pathname, handleMessage);
  }
});
