import ModalBase from './components/modal-base.vue';
import TheModals from './components/the-modals.vue';
import { useModalsStore } from './stores';

/**
 * Create public modals API.
 *
 * @param {Object} options
 * @param {Object} options.store
 * @param {Object} [options.dialogPropsMap]
 * @param {Object} [options.modal]
 * @return {Object}
 */
export const createModalsApi = ({ store, dialogPropsMap = {} }) => ({
  dialogPropsMap,

  isModalOpenedById: id => Boolean(store.byId[id]),

  show: modalPayload => store.show(modalPayload),

  hide: (payload = {}) => store.hide({ id: payload.id }),

  updateDialogProps: (payload = {}) =>
    store.updateDialogProps({
      id: payload.id,
      dialogProps: payload.dialogProps,
    }),

  updateModalConfig: (payload = {}) =>
    store.updateModalConfig({
      id: payload.id,
      config: payload.config,
    }),

  minimize: (payload = {}) => store.minimize({ id: payload.id }),

  maximize: (payload = {}) => store.maximize({ id: payload.id }),
});

/**
 * Modals plugin.
 */
export default {
  /**
   * Install modals plugin.
   *
   * @param {import('vue').App} app
   * @param {Object} options
   * @param {Object} [options.components]
   * @param {Object} [options.dialogPropsMap]
   * @param {string} [options.componentName]
   */
  install(app, { components = {}, dialogPropsMap = {}, componentName = 'the-modals' } = {}) {
    const modalsStore = useModalsStore();

    app.component(componentName, TheModals);
    app.component('modal-base', ModalBase);

    app.config.globalProperties.$modals = createModalsApi({
      store: modalsStore,
      dialogPropsMap,
    });

    app.provide('$modalsOptions', { dialogPropsMap });
    app.provide('$modalComponents', components);
  },
};
