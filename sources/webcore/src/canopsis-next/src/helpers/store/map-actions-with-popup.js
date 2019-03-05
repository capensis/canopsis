import { isFunction, isString } from 'lodash';

import mapActionsWith from './map-actions-with';

/**
 * Mapping actions with popups
 *
 * @param {Object} actions
 * @param {Object} messages
 */
export default function mapActionsWithPopup(actions, messages) {
  return mapActionsWith(actions, async function success(actionKey, result) {
    let popup = { text: this.$t('success.default') };

    if (messages[actionKey]) {
      let messageMethodWithContext;

      if (isFunction(messages[actionKey].success)) {
        messageMethodWithContext = messages[actionKey].success.bind(this);
      } else if (isFunction(messages[actionKey])) {
        messageMethodWithContext = messages[actionKey].bind(this);
      }

      const newPopup = messageMethodWithContext(result);

      if (isString(newPopup)) {
        popup.text = newPopup;
      } else {
        popup = newPopup;
      }
    }

    if (this.addSuccessPopup) {
      await this.addSuccessPopup(popup);
    }

    return result;
  }, async function error(actionKey, err) {
    let popup = { text: this.$t('errors.default') };

    if (messages[actionKey]) {
      let messageMethodWithContext;

      if (isFunction(messages[actionKey].error)) {
        messageMethodWithContext = messages[actionKey].error.bind(this);
      }

      const newPopup = messageMethodWithContext(err);

      if (isString(newPopup)) {
        popup.text = newPopup;
      } else {
        popup = newPopup;
      }
    }

    if (this.addErrorPopup) {
      await this.addErrorPopup(popup);
    }

    throw err;
  });
}

export function createMapActionsWithPopup(messages = {}) {
  return (actions, anotherMessages) => mapActionsWithPopup(actions, { ...messages, ...anotherMessages });
}
