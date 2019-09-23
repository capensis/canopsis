import { get } from 'lodash';

export default {
  install(Vue, { store }) {
    Object.defineProperty(Vue.prototype, '$modals', {
      get() {
        return {
          show: modal => store.dispatch('modal/show', modal),
          hide: ({ id } = {}) => store.dispatch('modal/hide', { id: id || get(this.modal, 'id') }),
        };
      },
    });
  },
};
