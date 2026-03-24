<template>
  <v-dialog
    v-model="isOpen"
    v-bind="dialogProps"
  >
    <!-- @slot use this slot default-->
    <slot />
  </v-dialog>
</template>

<script>
import { computed, inject, onMounted, ref } from 'vue';

import { DEFAULT_SIDEBAR_DRAWER_WIDTH } from '@/config';

import { useStore } from '@/hooks/store';
import { useModals } from '@/hooks/modals';
import { useSidebar } from '@/hooks/sidebar';

/**
 * Wrapper for each modal window
 *
 * @prop {Object} modal - The current modal object
 * @prop {Object} [dialogProps={}] - Properties for vuetify v-dialog
 */
export default {
  props: {
    modal: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const store = useStore();
    const modals = useModals();
    const sidebar = useSidebar();

    const clickOutside = inject('$clickOutside');

    const ready = ref(false);

    const modalSidebar = computed(() => store.getters[`${sidebar.moduleName}/sidebarsById`][props.modal.id]);

    const paddingRight = computed(() => (
      modalSidebar.value?.minimized || !modalSidebar.value?.config?.minimizable
        ? 0
        : modalSidebar.value?.config?.width || DEFAULT_SIDEBAR_DRAWER_WIDTH
    ));

    const isOpen = computed({
      get: () => !props.modal.hidden && ready.value,
      set: () => modals.hide({ id: props.modal.id }),
    });

    const dialogProps = computed(() => {
      const defaultDialogProps = {
        maxWidth: 700,
        attach: '.modals-wrapper',
        absolute: true,
        retainFocus: false,
      };
      const { dialogPropsMap = {} } = modals;
      const { name, dialogProps: modalDialogProps, minimized } = props.modal;

      const merged = {
        ...defaultDialogProps,
        ...dialogPropsMap[name],
        ...modalDialogProps,

        customCloseConditional: (...args) => clickOutside.call(...args),
      };

      if (paddingRight.value) {
        merged.contentWrapperStyle = {
          ...merged.contentWrapperStyle,

          paddingRight: `${paddingRight.value}px`,
        };
      }

      return {
        ...merged,

        hideOverlay: merged.hideOverlay || minimized,
        ignoreClickOutside: merged.ignoreClickOutside || minimized,
        contentWrapperClass: minimized ? 'v-dialog__content--minimized' : '',
      };
    });

    onMounted(() => ready.value = true);

    return {
      isOpen,
      dialogProps,
    };
  },
};
</script>

<style lang="scss">
.v-dialog .v-card__title {
  .headline {
    word-break: break-word;
  }
}
</style>
