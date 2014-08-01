describe("A suite", function() {
  it("contains spec with an expectation", function() {
    expect(true).toBe(true);
  });
});

describe("A WebSocket", function() {
  it("is defined", function() {
    expect(window.SOCKET).toBeDefined();
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
    spy = spyOnEvent('.ui-item', 'click');
    item = $('.ui-item').first();
    done();
  });

  it("show when an item is clicked", function(done) {
    setTimeout(function() {
      item.click();
      expect('click').toHaveBeenTriggeredOn('.ui-item');

      // var children = item.find('.ui-item-child');
      // expect(children.length).toBeGreaterThan(0);

      // expect(item.find('form input[name="text"]')).toBeFocused();
      done();
    }, 500);
  });

  it("hide when item is clicked again", function(done) {
    setTimeout(function() {
      item.click();
      expect('click').toHaveBeenTriggeredOn('.ui-item');
      done();
    }, 1000);
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
