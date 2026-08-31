import { computed, inject, unref } from 'vue';

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

    hide: (payload = {}) => store.hide({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    updateDialogProps: (payload = {}) => store.updateDialogProps({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    updateModalConfig: (payload = {}) => store.updateModalConfig({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    minimize: (payload = {}) => store.minimize({
        ...payload,
        id: getCurrentModalId(payload),
      }),

    maximize: (payload = {}) => store.maximize({
        ...payload,
        id: getCurrentModalId(payload),
      }),
  };
};

/**
 * Custom hook for managing inner modal states and actions within a Vue component.
 *
 * @param {Object} props - The properties passed to the inner modal, including the modal object.
 * @returns {Object} An object containing the modal's computed properties and a method to close the modal.
 */
export const useInnerModal = (props) => {
  const modals = useModals();

  const modal = computed(() => props.modal);
  const config = computed(() => modal.value?.config ?? {});

  const close = () => modals.hide({ id: props.modal?.id });

  return {
    modals,
    close,
    modal,
    config,
  };
};
