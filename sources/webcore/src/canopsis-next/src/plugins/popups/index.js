import { POPUP_TYPES } from '@/constants';

export default {
  install(Vue, { store }) {
    Object.defineProperty(Vue.prototype, '$popups', {
      get() {
        return {
          add: popup => store.dispatch('popup/add', popup),
          remove: ({ id } = {}) => store.dispatch('popup/remove', { id }),

          addSuccess: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.success }),
          addInfo: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.info }),
          addWarning: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.warning }),
          addError: popup => store.dispatch('popup/add', { ...popup, type: POPUP_TYPES.error }),
        };
      },
    });
  },
};
