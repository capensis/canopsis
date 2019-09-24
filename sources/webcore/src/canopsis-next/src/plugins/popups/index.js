import { POPUP_TYPES } from '@/constants';

export default {
  install(Vue, { store }) {
    Object.defineProperty(Vue.prototype, '$popups', {
      get() {
        return {
          add: popup => store.dispatch('popup/add', popup),
          remove: ({ id } = {}) => store.dispatch('popup/remove', { id }),

          success: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.success }),
          info: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.info }),
          warning: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.warning }),
          error: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.error }),
        };
      },
    });
  },
};
