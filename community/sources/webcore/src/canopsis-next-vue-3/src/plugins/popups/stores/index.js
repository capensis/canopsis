import { defineStore } from 'pinia';
import { ref } from 'vue';

import { POPUP_AUTO_CLOSE_DELAY } from '@/config';
import { POPUP_TYPES } from '@/constants';
import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} Popup
 * @property {string} id
 * @property {string} type
 * @property {string} [text]
 * @property {*} [context]
 * @property {number|boolean} [autoClose]
 */

/**
 * @typedef {Object} AddPopupPayload
 * @property {string} [id]
 * @property {string} type
 * @property {string} [text]
 * @property {*} [context]
 * @property {number|boolean} [autoClose]
 */

export const usePopupsStore = defineStore('popups', () => {
  const defaultCloseTimeByType = ref({
    [POPUP_TYPES.error]: POPUP_AUTO_CLOSE_DELAY,
    [POPUP_TYPES.success]: POPUP_AUTO_CLOSE_DELAY,
    [POPUP_TYPES.info]: POPUP_AUTO_CLOSE_DELAY,
    [POPUP_TYPES.warning]: POPUP_AUTO_CLOSE_DELAY,
  });

  /** @type {import('vue').Ref<Popup[]>} */
  const popups = ref([]);

  /**
   * Add popup.
   *
   * @param {AddPopupPayload} payload
   */
  const add = ({
    id = uid('popup'),
    type,
    text,
    context,
    autoClose = defaultCloseTimeByType.value[type],
  } = {}) => {
    popups.value.push({
      id,
      type,
      text,
      context,
      autoClose,
    });
  };

  /**
   * Remove popup by id.
   *
   * @param {{ id: string }} payload
   */
  const remove = ({ id }) => {
    popups.value = popups.value.filter(popup => popup.id !== id);
  };

  /**
   * Set default auto close time.
   *
   * @param {{ type: string, time?: number|boolean }} payload
   */
  const setDefaultCloseTime = ({ type, time = POPUP_AUTO_CLOSE_DELAY }) => {
    defaultCloseTimeByType.value[type] = time;
  };

  /**
   * Add success popup.
   *
   * @param {Omit<AddPopupPayload, 'type'>} popup
   */
  const success = (popup) => {
    add({ ...popup, type: POPUP_TYPES.success });
  };

  /**
   * Add info popup.
   *
   * @param {Omit<AddPopupPayload, 'type'>} popup
   */
  const info = (popup) => {
    add({ ...popup, type: POPUP_TYPES.info });
  };

  /**
   * Add warning popup.
   *
   * @param {Omit<AddPopupPayload, 'type'>} popup
   */
  const warning = (popup) => {
    add({ ...popup, type: POPUP_TYPES.warning });
  };

  /**
   * Add error popup.
   *
   * @param {Omit<AddPopupPayload, 'type'>} popup
   */
  const error = (popup) => {
    add({ ...popup, type: POPUP_TYPES.error });
  };

  return {
    defaultCloseTimeByType,
    popups,

    add,
    remove,
    setDefaultCloseTime,

    success,
    info,
    warning,
    error,
  };
});
