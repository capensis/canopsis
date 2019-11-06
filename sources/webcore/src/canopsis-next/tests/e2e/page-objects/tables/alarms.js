// https://nightwatchjs.org/guide/#working-with-page-objects

const el = require('../../helpers/el');

const commands = {
  clickAlarmListHeaderCell(header) {
    this.api
      .useXpath()
      .customClick(this.el('@alarmListHeaderCell', header));

    return this;
  },

  setAllCheckbox(checked) {
    return this.getAttribute('@selectAllCheckboxInput', 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick('@selectAllCheckbox');
      }
    });
  },

  setRowCheckbox(id, checked) {
    return this.getAttribute(this.el('@alarmListRowCheckboxInput', id), 'aria-checked', ({ value }) => {
      if (value !== String(checked)) {
        this.customClick(this.el('@alarmListRowCheckbox', id));
      }
    });
  },

  clickOnAlarmRow(id) {
    return this.customClick(this.el('@alarmListRow', id));
  },

  clickOnAlarmRowCell(id, column) {
    return this.customClick(this.el('@alarmListRow', id, column));
  },

  setItemPerPage(index) {
    return this.customClick('@itemsPerPage')
      .waitForElementVisible(this.el('@optionSelect', index))
      .customClick(this.el('@optionSelect', index));
  },

  clickOnMassAction(index) {
    return this.customClick(this.el('@massActionsPanelItem', index));
  },

  clickOnSharedAction(id, index) {
    return this.customClick(this.el('@rowActionsSharedPanelItem', id, index));
  },

  clickOnDropDownActions(id, index) {
    return this
      .customClick(this.el('@rowMoreActionsButton', id))
      .customClick(this.el('@rowDropDownActions', index));
  },

  el,
};

module.exports = {
  elements: {
    optionSelect: '.menuable__content__active .v-select-list [role="listitem"]:nth-of-type(%s)',

    alarmsListTable: sel('alarmsListTable'),

    alarmListHeaderCell: './/*[@data-test=\'alarmsListTable\']//thead//tr//th[@role=\'columnheader\']//span[contains(text(), \'%s\')]',
    selectAllCheckboxInput: `${sel('alarmsListTable')} thead tr th:first-of-type .v-input--selection-controls__input input`,
    selectAllCheckbox: `${sel('alarmsListTable')} thead tr th:first-of-type .v-input--selection-controls__input`,

    alarmListRow: `${sel('alarmListRow-%s')}`,
    alarmListRowCheckbox: `${sel('alarmListRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')}`,
    alarmListRowCheckboxInput: `${sel('alarmListRow-%s')} ${sel('rowCheckbox')} ${sel('vCheckboxFunctional')} input`,
    alarmListRowColumn: `${sel('alarmListRow-%s')} ${sel('alarmValue-%s')}`,

    massActionsPanelItem: `${sel('alarmsWidget')} ${sel('massActionsPanel')} ${sel('actionsPanelItem')}:nth-of-type(%s)`,
    rowActionsSharedPanelItem: `${sel('alarmListRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('actionsPanelItem')}:nth-of-type(%s)`,
    rowMoreActionsButton: `${sel('alarmListRow-%s')} ${sel('sharedActionsPanel')} .layout ${sel('dropDownActionsButton')}`,
    rowDropDownActions: `.menuable__content__active ${sel('dropDownActions')} ${sel('actionsPanelItem')}:nth-of-type(%s)`,
  },
  commands: [commands],
};
