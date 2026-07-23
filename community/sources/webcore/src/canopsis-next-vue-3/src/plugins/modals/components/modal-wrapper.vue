<script setup>
import { computed, inject, onMounted, ref } from 'vue';

import { useModals } from '@/hooks/modals';

const props = defineProps({
  modal: {
    type: Object,
    required: true,
  },
});

const clickOutside = inject('$clickOutside', null);
const { dialogPropsMap = {} } = inject('$modalsOptions', {});
const modals = useModals();

const ready = ref(false);
const hasClickOutsideHandlers = computed(() => Boolean(clickOutside?.handlersCount?.value));

const isOpen = computed({
  get: () => !props.modal.hidden && ready.value,
  set: (value) => {
    if (!value) {
      modals.hide({ id: props.modal.id });
    }
  },
});

const dialogProps = computed(() => {
  const defaultDialogProps = {
    maxWidth: 900,
    attach: '.modals-wrapper',
    absolute: true,
    retainFocus: false,
    scrollable: true,
  };

  const { name, dialogProps: modalDialogProps = {}, minimized } = props.modal;
  const { autoHeight, ...modalDialogPropsFromMap } = dialogPropsMap[name] ?? {};

  const merged = {
    ...defaultDialogProps,
    ...modalDialogPropsFromMap,
    ...modalDialogProps,
  };

  const contentClass = [
    merged.contentClass,
    autoHeight && 'v-dialog__content--auto-height',
    minimized && 'v-dialog__content--minimized',
  ]
    .filter(Boolean)
    .join(' ');

  return {
    ...merged,

    contentClass,
    scrim: minimized ? false : merged.scrim,
    persistent: minimized || merged.persistent || hasClickOutsideHandlers.value,
    noClickAnimation: hasClickOutsideHandlers.value || merged.noClickAnimation,
    closeOnBack: !minimized && merged.closeOnBack,
  };
});

/**
 * Handles outside click.
 *
 * @param {PointerEvent} event
 */
const handleClickOutside = (event) => {
  if (props.modal.minimized) {
    return;
  }

  if (clickOutside && !clickOutside.call(event)) {
    return;
  }

  modals.hide({ id: props.modal.id });
};

onMounted(() => (ready.value = true));
</script>

<template>
  <v-dialog v-model="isOpen" v-bind="dialogProps" @click:outside="handleClickOutside">
    <slot />
  </v-dialog>
</template>

<style lang="scss">
.modals-wrapper .v-dialog {
  &:has(> .v-dialog__content--minimized) {
    align-items: flex-end;
    justify-content: center;
    pointer-events: none;
  }

  > .v-overlay__content.v-dialog__content--minimized {
    margin: 8px 8px 0 8px;
    position: relative;
    height: auto;
    max-height: none;
    max-width: 360px;
    width: fit-content;
    min-width: 0;
    pointer-events: all;
    z-index: inherit !important;
    padding: 0 !important;
    transform: none !important;
    box-shadow: none;
    transition: none;

    > .v-form,
    > .v-card {
      height: auto;
      min-height: auto !important;
      max-height: none;
      max-width: 360px;
      width: fit-content;
      min-width: 0;
      flex: 0 0 auto;
      overflow: hidden;
    }

    .v-card {
      min-height: auto !important;
      max-width: 360px;
      width: fit-content;

      &.fill-min-height {
        min-height: auto;
      }
    }

    .v-card-title,
    .v-card__title {
      display: flex;
      align-items: center;
      min-height: 40px;
      padding: 0 10px;
      width: fit-content;
      max-width: 360px;

      .modal-wrapper__title,
      .modal-title-buttons,
      .modal-title-button__wrapper {
        display: flex;
        align-items: center;
      }

      .text-h5,
      .headline {
        font-size: 16px !important;
        line-height: 1.25;
      }
    }
  }

  &.v-dialog--scrollable > .v-overlay__content > .v-card,
  &.v-dialog--scrollable > .v-overlay__content > .v-form > .v-card {
    overflow: hidden;
  }

  > .v-overlay__content.v-dialog__content--auto-height {
    height: auto !important;
    max-height: calc(100% - 48px);
    overflow: hidden;
    padding: 0;
    display: flex;
    flex-direction: column;
    min-height: 0;

    > .v-form {
      height: auto;
      max-height: 100%;
      min-height: 0;
      display: flex;
      flex-direction: column;

      > .v-card {
        min-height: unset !important;
        height: auto;
        max-height: 100%;
        flex: 1 1 auto;
        min-height: 0;
        overflow: hidden;
      }
    }

    > .v-card {
      min-height: unset !important;
      height: auto;
      max-height: 100%;
      flex: 1 1 auto;
      min-height: 0;
      overflow: hidden;
    }
  }

  > .v-overlay__content {
    overflow: hidden;
    display: flex;
    flex-direction: column;
    min-height: 0;

    > .v-form {
      display: flex;
      flex-direction: column;
      height: 100%;
      min-height: 0;

      > .v-card {
        flex: 1 1 0;
        min-height: 0;
        overflow: hidden;
      }
    }

    > .v-card {
      flex: 1 1 0;
      min-height: 0;
      overflow: hidden;
    }

    .v-card-title,
    .v-card__title,
    .v-card-actions {
      flex-shrink: 0;
    }

    .v-card-text {
      flex: 1 1 auto;
      min-height: 0;
      overflow-y: auto;
    }

    .headline {
      word-break: break-word;
    }
  }
}
</style>
