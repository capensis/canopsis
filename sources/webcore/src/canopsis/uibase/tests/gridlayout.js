module('uibase', {
  beforeEach: function() {
    $('.modal-backdrop').remove();
    visit('/userview/view.event');
  },
  afterEach: function() {
    $('.modal-backdrop').remove();
  }
});

test('Grid layout children options', function() {
    expect(2);

    createNewView('gridlayout_childrenpackingoptions');

    activateEditMode();
    changeEditorForKey('columnXS', 'defaultpropertyeditor').then(function() {
      changeEditorForKey('columnMD', 'defaultpropertyeditor').then(function() {
        changeEditorForKey('columnLG', 'defaultpropertyeditor').then(function() {
          click('.btn-add-widget');
          waitForElement('.form .ember-text-field').then(function(){
            fillIn('.form .ember-text-field', 'widgetcontainer');
            click('.form .panel-default a');
            click('.form .list-group-item a');
            click('.form #mixins_tab a');
            fillIn('.tab-pane.active input', 'gridlayout');
            click('.form .panel-default .panel-heading a');
            click('.form .panel-default .list-group a');
            click('.form .btn-submit');

            waitForElementRemoval('.modal-backdrop').then(function() {
              click('.widgetslot .btn-add-widget');
              waitForElement('.form .ember-text-field').then(function(){
                fillIn('.form .ember-text-field', 'text');
                click('.form .panel-default a');
                click('.form .list-group-item a');
                click('.form .btn-submit');

                waitForElementRemoval('.modal-backdrop').then(function() {
                  click('.widgetslot .widgetslot:first .box-header .dropdown-toggle');
                  click('.widgetslot .widgetslot:first .box-header .dropdown-menu a');
                  waitForElement('.form .ember-text-field').then(function(){
                    fillIn('.form input[name=columnXS]', 4);
                    fillIn('.form input[name=columnMD]', 4);
                    fillIn('.form input[name=columnLG]', 4);
                    click('.form .btn-submit');

                    waitForElementRemoval('.modal-backdrop').then(function() {
                      click('.widgetslot .widgetslot:first .fa-files-o');
                      waitMilliseconds(1000).then(function() {
                        click('.widgetslot .widgetslot:last .box-header .dropdown-toggle');
                        click('.widgetslot .widgetslot:last .box-header .dropdown-menu a');
                        waitForElement('.form .ember-text-field').then(function(){

                          fillIn('.form input[name=columnXS]', 6);
                          fillIn('.form input[name=columnMD]', 6);
                          fillIn('.form input[name=columnLG]', 6);

                          click('.form .btn-submit');
                          waitMilliseconds(1000).then(function() {
                            equal(find('.col-md-4.col-xs-4.col-lg-4').length, 1, 'there is one width with packing widths set to [4,4,4]');
                            equal(find('.col-md-6.col-xs-6.col-lg-6').length, 1, 'there is one width with packing widths set to [6,6,6]');
                          });
                        });
                      });
                    });
                  });
                });
              });
            });
          });
        });
      });
    });
});
