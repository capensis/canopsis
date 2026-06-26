import { usePopupsStore } from './stores';
import ThePopups from './components/the-popups.vue';

/**
 * Creates public popups API.
 *
 * @param {ReturnType<typeof usePopupsStore>} store
 * @return {Object}
 */
export const createPopupsApi = store => ({
  add: popup => store.add(popup),
  remove: popup => store.remove(popup),
  success: popup => store.success(popup),
  info: popup => store.info(popup),
  warning: popup => store.warning(popup),
  error: popup => store.error(popup),
  setDefaultCloseTime: popup => store.setDefaultCloseTime(popup),
});

/**
 * Popups plugin.
 */
export default {
  /**
   * Install popups plugin.
   *
   * @param {import('vue').App} app
   * @param {Object} options
   * @param {string} [options.componentName]
   */
  install(app, { componentName = 'the-popups' } = {}) {
    const popupsStore = usePopupsStore();
    const popupsApi = createPopupsApi(popupsStore);

    app.component(componentName, ThePopups);

    app.config.globalProperties.$popups = popupsApi;

    app.provide('$popups', popupsApi);
  },
};
