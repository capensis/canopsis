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

import { useModals } from '@/hooks/modals';

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
    const modals = useModals();

    const clickOutside = inject('$clickOutside');

    const ready = ref(false);

    const isOpen = computed({
      get: () => !props.modal.hidden && ready.value,
      set: () => modals.hide({ id: props.modal.id }),
    });

    const dialogProps = computed(() => {
      const defaultDialogProps = {
        maxWidth: 900,
        attach: '.modals-wrapper',
        absolute: true,
        retainFocus: false,
      };
      const { dialogPropsMap = {} } = modals;
      const { name, dialogProps: modalDialogProps, minimized } = props.modal;
      const { autoHeight, ...modalDialogPropsFromMap } = dialogPropsMap[name] ?? {};

      const merged = {
        ...defaultDialogProps,
        ...modalDialogPropsFromMap,
        ...modalDialogProps,

        customCloseConditional: (...args) => clickOutside.call(...args),
      };
      if (autoHeight) {
        merged.contentWrapperClass = `${merged.class ?? ''} v-dialog__content--auto-height`.trim();
      }

      if (minimized) {
        merged.contentWrapperClass = `${merged.contentWrapperClass ?? ''} v-dialog__content--minimized`.trim();
      }

      return {
        ...merged,

        hideOverlay: merged.hideOverlay || minimized,
        ignoreClickOutside: merged.ignoreClickOutside || minimized,
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
.v-dialog__content--auto-height .v-dialog {
  height: auto !important;

  > .v-form {
    height: auto;

    > .v-card {
      min-height: unset !important;
      height: auto;
      flex: 0 0 auto;
    }
  }

  > .v-card {
    min-height: unset !important;
    height: auto;
    flex: 0 0 auto;
  }
}

.v-dialog {
  height: 70vh;

  &:not(.v-dialog--auto-height) > {
    .v-form, .v-card {
      min-height: 100%;
    }
  }

  & > .v-form {
    display: flex;
    flex-direction: column;
    height: 100%;

    & > .v-card {
      flex: 1 1 0;
      min-height: 0;
    }
  }

  .v-card__title {
    .headline {
      word-break: break-word;
    }
  }
}
</style>
