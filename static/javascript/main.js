var isSelectMode = false;
var tabFocusedItemHash = null;

var handleMessage = function(data) {
  var item = $('#'+data.item.hash);
  if (item.length > 0) {
    ItemManager.updateChildren(data.items, item);
  } else {
    var root = $('.ui-root-items');
    ItemManager.updateItems(data.items, root);
  }
};

var isLight = function(text) {
  if (text == "") {
    text = "110,171,221";
  }
  var rgb = text.split(',');
  var hsl = rgbToHsl(rgb[0], rgb[1], rgb[2]);
  return hsl[2] > 0.75;
};

function rgbToHsl(r, g, b){
  r /= 255, g /= 255, b /= 255;
  var max = Math.max(r,g,b);
  var min = Math.min(r,g,b);
  var h, s, l = (max + min) / 2;

  if (max == min) {
    h = s = 0; // achromatic
  } else {
    var d = max - min;
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min);
    switch(max){
      case r: h = (g - b) / d + (g < b ? 6 : 0); break;
      case g: h = (b - r) / d + 2; break;
      case b: h = (r - g) / d + 4; break;
    }
    h /= 6;
  }

  return [h,s,l];
};

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
};

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
};

// UpdateChildren checks all the child items on screen and adds any missing
// children from the current dataset.
ItemManager.updateChildren = function(items, parent) {
  parent.find('.ui-item').remove();
  var anchor = parent.find('> a');
  for (var i=0; i<items.length; i++) {
    var item = Item.html(items[i], 'ui-item-child');
    anchor.after(item);
  }
};

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
};

// Color saves an assigned color to an item.
ItemManager.color = function(e) {
  e.preventDefault();

  var target = $(this);
  var item = $('#'+e.data.hash);
  var color = $(e.target).data('color');

  item.find('> a').css('background-color', 'rgba('+color+',1)');
  if (isLight(color)) {
    item.addClass('ui-item-light');
  } else {
    item.removeClass('ui-item-light');
  }

  e.data.color = color;
  var data = $.param(e.data, true);
  Item.save(data);
};

ItemManager.handleContextMenu = function(e) {
  e.preventDefault();

  var item = $(e.target).parent();
  if (!item.hasClass('ui-item')) {
    return;
  }

  Menu.show(item.data(), {
    'Colors': ItemManager.color,
    'View': Item.view,
    'Edit': Item.edit,
    'Delete': Item.delete
  }, e);
};

ItemManager.handleQuestionMenu = function(e) {
  e.preventDefault();
  e.stopPropagation();
  var item = $(e.target).parents('.ui-item');
  Menu.show(item.data(), {
    'ToField': null,
    'MessageField': null,
  }, e);
};

var Item = {};

// HandleClick kicks off an xhr request to pull down child items, if they
// exist. Or hides items if they're already showing.
Item.handleClick = function(e) {
  e.preventDefault();
  e.stopPropagation();

  Menu.clear();

  if (e.target.tagName == 'DIV' || e.target.tagName == 'INPUT') {
    return;
  }

  Item.Focus($(this));
};

Item.Focus = function(item) {
  var url = '/i/'+item.data('hash');

  if (isSelectMode) {
    item.toggleClass('ui-item-selected');
    return;
  }

  if (item.hasClass('ui-item-expanded')) {
    item.find('.ui-item').remove();
    item.find('.ui-item-form').remove();
    item.removeClass('ui-item-expanded');
    window.SOCKET.unsubscribe(url);
  } else {
    $.ajax({
      'url': url,
      'success': ItemManager.addItems
    });
    window.SOCKET.subscribe(url, handleMessage);
  }
}

// HTML returns HTML necessary to render an item.
Item.html = function(data, extraClass) {
  var text = data.text;
  var isQuestion = false;
  if (data.text.slice(-1) == "?") {
    text = text.slice(0,-1);
    isQuestion = true;
  }

  if (isLight(data.color)) {
    extraClass += ' ui-item-light'
  }

  var item = $(''+
    '<div class="ui-item '+extraClass+'" id="'+data.hash+'">'+
      '<a href="/i/'+data.hash+'" style="background-color:rgba('+data.color+',.75); border-color:rgba('+data.color+',1);">'+text+'</a>'+
    '</div>');
  item.data({
    'hash': data.hash,
    'parent': data.parent,
    'root': data.root,
    'user': data.user,
    'text': data.text,
    'color': data.color
  });

  if (isQuestion) {
    item.find('a').append('<span class="ui-item-question">?</span>');
  }

  return item;
};

// Save saves a new item.
Item.save = function(data, complete) {
  $.post('/i/save', data, function(data) {
    if (complete) {
      complete(data);
    }

    window.SOCKET.request('/i/'+data.item.parent);
  }.bind(this));
};

// Edit edits an item.
Item.edit = function(e) {
  e.preventDefault();
  alert("Not implemented yet :(");
};

// Deletes removes an item.
Item.delete = function(e) {
  e.preventDefault();
  $.post('/i/'+e.data.hash+'/delete', function(data) {
    if (data.error) {
      console.log(data.error);
    } else {
      var item = $('#'+e.data.hash);
      item.remove();
      window.SOCKET.request('/i/'+e.data.parent);
    }
  });
};

Item.view = function(e) {
  e.preventDefault();
  window.location = '/i/'+e.data.hash;
};

// Insert adds an item into a given parent node.
Item.insert = function(item, parent) {
  parent.prepend(item);
};

var MessageManager = {};

MessageManager.handleContextMenu = function(e) {
  e.preventDefault();
  Menu.show($(this).data(), {
    'Edit': Message.edit,
    'Delete': Message.delete
  }, e);
};

var Message = {};

Message.edit = function(e) {
  e.preventDefault();
  alert("Not implemented yet :(");
};

Message.delete = function(e) {
  e.preventDefault();
  $.post('/m/'+e.data.hash+'/delete', function(data) {
    if (data.error) {
      console.log(data.error);
    } else {
      var message = $('#'+e.data.hash);
      message.remove();
      // window.SOCKET.request('/i/'+target.data('parent'));
    }
  });
};

var Menu = {};

Menu.show = function(data, options, e) {
  Menu.clear();

  var menu = $('<div class="ui-menu"><ul></ul></div>');
  var ul = menu.find('ul');

  for (option in options) {
    var item = $('<li><a class="ui-menu-'+option.toLowerCase()+'" href="#">'+option+'</a></li>');
    item.on('click', data, options[option]);

    if (option == 'Colors') {
      item.find('a').text('');
      item.find('a').append(''+
        '<span data-color="79,89,119" class="ui-color ui-color-dark">Dark</span>'+
        '<span data-color="110,171,221" class="ui-color ui-color-blue ui-color-selected">Blue</span>'+
        '<span data-color="152,225,178" class="ui-color ui-color-green">Green</span>'+
        '<span data-color="255,231,159" class="ui-color ui-color-yellow">Yellow</span>'+
        '<span data-color="246,155,129" class="ui-color ui-color-red">Red</span>'
        );
    }
    ul.append(item);
  }

  $('body').append(menu);

  if ((e.pageY + menu.height()) > window.innerHeight) {
    menu.css({'top': e.pageY - menu.height(), 'left': e.pageX});
  } else {
    menu.css({'top': e.pageY, 'left': e.pageX});
  }
};

Menu.clear = function() {
  $('.ui-menu').remove();
};

// HACK
$(function() {
  var body = $('body');

  // Items
  body.on('click', '.ui-item', Item.handleClick);
  body.on('submit', '.ui-item-form', ItemManager.submit);
  body.on('contextmenu', '.ui-item', ItemManager.handleContextMenu);

  // Messages
  body.on('contextmenu', '.ui-message', MessageManager.handleContextMenu);

  // Menus
  $(document).on('click', Menu.clear);

  // Input fields
  $('.ux-focus').focus();

  // Keyboard shortcuts
  body.on({
    'keydown': function(e) {
      if (e.keyCode == 16) { // Shift
        console.log('ON');
        isSelectMode = true;
      }
      if (e.keyCode == 9) { // Tab
        e.preventDefault();
        var item;
        if (tabFocusedItemHash) {
          var previousItem = $('#'+tabFocusedItemHash);
          item = previousItem.next();
          Item.Focus(previousItem);
        } else {
          item = $('.ui-item').first();
        }
        tabFocusedItemHash = item.data('hash');
        Item.Focus(item);
      }
    },
    'keyup': function(e) {
      if (e.keyCode == 16) { // Shift
        console.log('OFF');
        isSelectMode = false;
      }
    }
  });

  // WebSocket
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
