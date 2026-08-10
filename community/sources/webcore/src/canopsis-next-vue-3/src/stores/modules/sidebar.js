import { defineStore } from 'pinia';
import { computed, ref } from 'vue';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { uid } from '@/helpers/uid';

/**
 * @typedef {Object} Sidebar
 * @property {string} id
 * @property {string} name
 * @property {Object} config
 * @property {boolean} hidden
 */

export const useSidebarStore = defineStore('sidebar', () => {
  const allIds = ref([]);
  const byId = ref({});

  const sidebars = computed(() => allIds.value.map(id => byId.value[id]));

  /**
   * Show a sidebar by component name with a stable id.
   *
   * @param {Object} [payload={}]
   * @param {string} payload.name
   * @param {Object} [payload.config={}]
   * @param {string} [payload.id=uid()]
   */
  const show = ({ id = uid('sidebar'), name, config = {} } = {}) => {
    const openedSidebar = sidebars.value.find(sidebar => sidebar && !sidebar.hidden);

    if (openedSidebar) {
      return openedSidebar.id;
    }

    allIds.value
      .filter(sidebarId => byId.value[sidebarId]?.hidden)
      .forEach(sidebarId => hideCompleted({ id: sidebarId }));

    byId.value[id] = {
      id,
      name,
      config,
      hidden: false,
    };

    allIds.value.push(id);

    return id;
  };

  /**
   * Hide sidebar by id.
   *
   * @param {Object} [payload={}]
   * @param {string} payload.id
   */
  const hide = ({ id } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    if (!byId.value[id]) {
      return;
    }

    byId.value[id].hidden = true;

    setTimeout(() => {
      if (byId.value[id]) {
        hideCompleted({ id });
      }
    }, VUETIFY_ANIMATION_DELAY);
  };

  /**
   * Remove hidden sidebar from store.
   *
   * @param {Object} payload
   * @param {string} payload.id
   */
  const hideCompleted = ({ id }) => {
    allIds.value = allIds.value.filter(value => value !== id);

    delete byId.value[id];
  };

  /**
   * Update sidebar config.
   *
   * @param {Object} [payload={}]
   * @param {string} payload.id
   * @param {Object} [payload.config={}]
   */
  const updateConfig = ({ id, config = {} } = {}) => {
    if (!id) {
      throw new Error('Missed required parameter');
    }

    if (!byId.value[id]) {
      return;
    }

    byId.value[id].config = {
      ...byId.value[id].config,
      ...config,
    };
  };

  return {
    allIds,
    byId,
    sidebars,

    show,
    hide,
    hideCompleted,
    updateConfig,
  };
});
