import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';

import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} Modal
 * @property {string} id
 * @property {string} name
 * @property {Function} [onHide]
 * @property {Object} config
 * @property {Object} dialogProps
 * @property {boolean} hidden
 * @property {boolean} minimized
 */

/**
 * @typedef {Object} ShowModalPayload
 * @property {string} name
 * @property {Function} [onHide]
 * @property {Object} [config]
 * @property {Object} [dialogProps]
 * @property {string} [id]
 */

export const useModalsStore = defineStore('modals', () => {
  const allIds = ref([]);
  const byId = ref({});

  const modals = computed(() => allIds.value.map(id => byId.value[id]));
  const hasMaximizedModal = computed(() => allIds.value.some(id => !byId.value[id].minimized));

  /**
   * Show modal window by name with unique id.
   *
   * @param {ShowModalPayload} payload
   */
  const show = ({ name, onHide, config = {}, dialogProps = {}, id = uid('modal') } = {}) => {
    if (byId.value[id]) {
      maximize({ id });
      return;
    }

    byId.value[id] = {
      id,
      name,
      onHide,
      config,
      dialogProps,
      hidden: false,
      minimized: false,
    };

    allIds.value.push(id);
  };

  /**
   * Update dialog props.
   *
   * @param {{ id: string, dialogProps: Object }} payload
   */
  const updateDialogProps = ({ id, dialogProps } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    byId.value[id].dialogProps = dialogProps;
  };

  /**
   * Update modal config.
   *
   * @param {{ id: string, config?: Object }} payload
   */
  const updateModalConfig = ({ id, config = {} } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    if (!byId.value[id]) {
      return;
    }

    byId.value[id].config = config;
  };

  /**
   * Hide modal by id.
   *
   * @param {{ id: string }} payload
   */
  const hide = async ({ id } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    const modal = byId.value[id];

    if (!modal) {
      return;
    }

    modal.hidden = true;

    await modal.onHide?.();

    setTimeout(() => {
      if (byId.value[id]) {
        hideCompleted({ id });
      }
    }, VUETIFY_ANIMATION_DELAY);
  };

  /**
   * Complete modal hiding and remove it from store.
   *
   * @param {{ id: string }} payload
   */
  const hideCompleted = ({ id }) => {
    allIds.value = allIds.value.filter(value => value !== id);

    delete byId.value[id];
  };

  /**
   * Minimize modal by id.
   *
   * @param {{ id: string }} payload
   */
  const minimize = ({ id } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    byId.value[id].minimized = true;
  };

  /**
   * Maximize modal by id.
   *
   * @param {{ id: string }} payload
   */
  const maximize = ({ id } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    byId.value[id].minimized = false;
  };

  return {
    allIds,
    byId,

    modals,
    hasMaximizedModal,

    show,
    hide,
    hideCompleted,
    updateDialogProps,
    updateModalConfig,
    minimize,
    maximize,
  };
});
