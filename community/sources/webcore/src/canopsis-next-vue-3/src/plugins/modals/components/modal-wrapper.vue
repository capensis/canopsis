<template>
  <v-dialog v-model="isOpen" v-bind="dialogProps" @click:outside="handleClickOutside">
    <slot />
  </v-dialog>
</template>

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
  const {
    contentWrapperStyle,
    contentProps = {},
    ...dialogPropsWithoutContent
  } = merged;
  const preparedContentProps = contentProps ?? {};
  const rightOffset = contentWrapperStyle?.paddingRight;
  const hasContentRightOffset = rightOffset !== undefined && rightOffset !== null;
  const contentStyle = {
    ...preparedContentProps.style,
    ...contentWrapperStyle,
  };
  contentStyle.transition = [
    contentStyle.transition,
    'width 0.2s ease-in-out, margin-right 0.2s ease-in-out, max-width 0.2s ease-in-out',
  ].filter(Boolean).join(', ');

  if (hasContentRightOffset) {
    const normalizedRightOffset = rightOffset || 0;

    contentStyle.paddingRight = undefined;
    contentStyle.marginRight = `calc(${normalizedRightOffset} + 24px)`;
    contentStyle.width = `calc(100% - ${normalizedRightOffset} - 48px)`;
  }

  const contentClass = [
    dialogPropsWithoutContent.contentClass,
    autoHeight && 'v-dialog__content--auto-height',
    minimized && 'v-dialog__content--minimized',
  ]
    .filter(Boolean)
    .join(' ');

  return {
    ...dialogPropsWithoutContent,

    contentClass,
    contentProps: {
      ...preparedContentProps,
      style: contentStyle,
    },
    scrim: minimized ? false : dialogPropsWithoutContent.scrim,
    persistent: minimized || dialogPropsWithoutContent.persistent || hasClickOutsideHandlers.value,
    noClickAnimation: hasClickOutsideHandlers.value || dialogPropsWithoutContent.noClickAnimation,
    closeOnBack: !minimized && dialogPropsWithoutContent.closeOnBack,
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
      overflow-x: hidden;
    }

    .v-window,
    .v-window__container {
      max-width: 100%;
      min-width: 0;
      overflow: visible;
    }

    .v-window-item {
      max-width: 100%;
      min-width: 0;
      overflow: visible;
    }

    .v-window__container {
      height: auto !important;
    }

    .v-window__container--is-active {
      height: auto !important;
    }

    .v-window:has(
      .v-window-x-transition-enter-active,
      .v-window-x-transition-leave-active,
      .v-window-x-reverse-transition-enter-active,
      .v-window-x-reverse-transition-leave-active,
      .v-window-y-transition-enter-active,
      .v-window-y-transition-leave-active,
      .v-window-y-reverse-transition-enter-active,
      .v-window-y-reverse-transition-leave-active
    ) {
      overflow: hidden;
    }

    .headline {
      word-break: break-word;
    }
  }
}
</style>
