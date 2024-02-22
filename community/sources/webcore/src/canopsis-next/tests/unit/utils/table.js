/**
 * Select an advanced table
 *
 * @param {Wrapper} wrapper
 * @return {Wrapper}
 */
export const selectTable = wrapper => wrapper.find('.c-advanced-data-table, .v-datatable');

/**
 * Select an advanced table rows
 *
 * @param {Wrapper} wrapper
 * @return {WrapperArray}
 */
export const selectTableRows = wrapper => selectTable(wrapper)
  .findAll('tbody > tr:not(.v-data-table__expanded)');

/**
 * Select an advanced table row by index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectTableRowByIndex = (wrapper, index) => selectTableRows(wrapper)
  .at(index);

/**
 * Select an advanced table row select checkbox by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowCheckboxByIndex = (wrapper, index) => selectTableRowByIndex(wrapper, index)
  .find('td:first-child .v-data-table__checkbox');

/**
 * Select an advanced table row select expand button by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowExpandButtonByIndex = (wrapper, index) => selectTableRowByIndex(wrapper, index)
  .find('td .v-data-table__expand-icon');

/**
 * Select an advanced table mass actions
 *
 * @param {Wrapper} wrapper
 * @return {Wrapper}
 */
export const selectMassActions = wrapper => selectTable(wrapper)
  .find('.layout');

/**
 * Select an advanced table mass remove action
 *
 * @param {Wrapper} wrapper
 * @return {Wrapper}
 */
export const selectMassRemoveButton = wrapper => selectMassActions(wrapper)
  .find('c-action-btn-stub');

/**
 * Select an advanced table row actions by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowActionsByIndex = (wrapper, index) => selectTableRowByIndex(wrapper, index)
  .find('td:last-child');

/**
 * Select an advanced table row actions by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowActionsButtonsByIndex = (wrapper, index) => selectRowActionsByIndex(wrapper, index)
  .findAll('c-action-btn-stub');

/**
 * Select an advanced table row edit button by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowEditButtonByIndex = (wrapper, index) => selectRowActionsByIndex(wrapper, index)
  .find('c-action-btn-stub[type="edit"]');

/**
 * Select an advanced table row remove button by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowRemoveButtonByIndex = (wrapper, index) => selectRowActionsByIndex(wrapper, index)
  .find('c-action-btn-stub[type="delete"]');

/**
 * Select an advanced table row remove duplicate by row index
 *
 * @param {Wrapper} wrapper
 * @param {number} index
 * @return {Wrapper}
 */
export const selectRowDuplicateButtonByIndex = (wrapper, index) => selectRowActionsByIndex(wrapper, index)
  .find('c-action-btn-stub[type="duplicate"]');
