describe("A suite", function() {
  it("contains spec with an expectation", function() {
    expect(true).toBe(true);
  });
});

describe("A WebSocket", function() {
  beforeEach(function(done) {
    spyOn(window.SOCKET, 'onopen');
    setTimeout(function() {
      done();
    }, 1000);
  });

  it("is defined", function() {
    expect(window.SOCKET).toBeDefined();
  });

  it("is opened", function() {
    expect(window.SOCKET.onopen).toHaveBeenCalled();
  });
});

describe("An item page", function() {
  it("on load has a focused field", function() {
    expect($('#root-form input[name="text"]')).toBeFocused();
  });

  it("items have an id property", function() {
    expect($('.ui-item')).toHaveProp('id');
  });

  it("items have data", function() {
    expect($('.ui-item')).toHaveData('hash');
    expect($('.ui-item')).toHaveData('parent');
    expect($('.ui-item')).toHaveData('root');
    expect($('.ui-item')).toHaveData('text');
    expect($('.ui-item')).toHaveData('color');
  });
});

describe("Child items", function() {
  var spy = null;
  var item = null

  beforeEach(function(done) {
    item = $('.ui-item a').first();
    spyOn(item, 'click');

    setTimeout(function() {
      item.click();
      done();
    }, 1000);
  });

  it("show when an item is clicked", function(done) {
    expect(item.click).toHaveBeenCalled();
    console.log(item.parent().find('form'));
    done();
  });

  it("hide when item is clicked again", function(done) {
    expect(item.click).toHaveBeenCalled();
    done();
  });
});

// describe("Contextual menus", function() {
//   it("shows when items are right clicked", function() {
//     var spy = spyOnEvent('.ui-item', 'contextmenu');
//     var item = $('ui-item').first();
//     item.contextmenu();
//
//     expect($('#item-menu')).toBeVisible();
//   });
// });
