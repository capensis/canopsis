import { inject, unref } from 'vue';

import { useModalsStore } from '@/plugins/modals/stores';

export const useModals = () => {
  const store = useModalsStore();
  const modalContext = inject('$modals', null);

  const getCurrentModalId = (payload) => {
    if (payload?.id) {
      return payload.id;
    }

    return unref(modalContext)?.id;
  };

  return {
    dialogPropsMap: {},

    isModalOpenedById: id => Boolean(store.byId[id]),

    show: modal => store.show(modal),

    hide: (payload = {}) =>
      store.hide({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    updateDialogProps: (payload = {}) =>
      store.updateDialogProps({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    updateModalConfig: (payload = {}) =>
      store.updateModalConfig({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    minimize: (payload = {}) =>
      store.minimize({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    maximize: (payload = {}) =>
      store.maximize({
        ...payload,
        id: getCurrentModalId(payload),
      }),
  };
};
