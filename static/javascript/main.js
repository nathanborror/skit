
function handleMessage(data) {
  var parent = $('#'+data.skit.hash);
  replaceItems(data.children, parent);
}

// Kicks off an xhr request to pull down child items, if they exist.
// Or hides children if they're already showing.

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
  var parent = $('#'+data.skit.hash);
  var anchor = parent.find('> a');
  parent.addClass('ui-item-expanded');

  for (var i=0; i<data.children.length; i++) {
    var item = getItemHTML(data.children[i], 'ui-item-child');
    anchor.after(item);
  }

  // Copy form to items
  var form = $('#root-form').clone().appendTo(parent);
  form.attr('id', '');
  var parent_field = form.find('[name="parent"]');
  var root_field = form.find('[name="root"]');
  parent_field.val(data.skit.hash)
  root_field.val(data.skit.root)
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
  var item = $('<div class="ui-item '+extraClass+'" id="'+data.hash+'"><a href="/s/'+data.hash+'">'+data.text+'</a></div>');
  item.data({'hash': data.hash, 'parent': data.parent, 'root': data.root});
  return item;
}

// Saves a new item and inserts it into the DOM

function handleSave(e) {
  e.preventDefault();
  var form = $(this);

  $.post('/s/save', $(this).serialize(), function(data) {
    if (form.parent().prop('tagName') == 'ARTICLE') {
      var parent = $('.ui-root-items')
      var item = getItemHTML(data.skit);
      parent.append(item);
    } else {
      var parent = $('#'+data.skit.parent);
      var item = getItemHTML(data.skit, 'ui-item-child');
      form.before(item);
    }
    window.SOCKET.request('/s/'+data.skit.parent);
  }.bind(this));

  form.find('input[name="text"]').val("");
}

function handleDelete(e) {
  e.preventDefault();

  var target = $(this);
  $.post('/s/'+target.data('hash')+'/delete', function(data) {
    var item = $('#'+target.data('hash'));
    item.remove();
    window.SOCKET.request('/s/'+target.data('parent'));
  });
}

// Shows or hides a contextual menu

function handleContextMenu(e) {
  e.preventDefault();
  var menu = $('#menu');
  var item = $(this).parent();

  menu.find('.ui-item-share').attr('href', '/s/'+item.attr('id'));
  menu.find('.ui-item-share').data({'hash': item.data('hash'), 'parent': item.data('parent'), 'root': item.data('root')});

  menu.find('.ui-item-delete').attr('href', '/s/'+item.attr('id')+'/delete');
  menu.find('.ui-item-delete').data({'hash': item.data('hash'), 'parent': item.data('parent'), 'root': item.data('root')});

  menu.find('.ui-item-edit').attr('href', '/s/'+item.attr('id')+'/edit');
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
  $('body').on('click', '.ui-item a', handleItemClick);
  $('body').on('contextmenu', '.ui-item a', handleContextMenu);
  $('body').on('submit', '.ui-item-form', handleSave);
  $('body').on('click', '.ui-item-delete', handleDelete);

  $(document).on('click', function() {
    $('#menu').hide();
  });

  $('footer a').each(function() {
    var user = $(this);
    var color = handleColor(user.data('hash'));
    user.css('background', color);
  });

  $('.ux-focus').focus();


});
