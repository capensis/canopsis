<template>
  <v-card :class="wrapperClasses" class="modal-wrapper">
    <v-card-title
      v-if="$slots.title"
      :style="titleStyle"
      :class="titleClass ?? 'white--text'"
    >
      <div class="modal-wrapper__title text-h5">
        <div>
          <slot name="title" />
        </div>
        <div>
          <modal-title-buttons
            :minimize="minimize"
            :close="close"
          />
        </div>
      </div>
      <modal-mass-actions-panel :id="$modal.id" />
    </v-card-title>
    <template v-if="!$modal.minimized">
      <v-card-text
        v-if="$slots.text"
        key="text"
        :class="textClass"
      >
        <slot name="text" />
      </v-card-text>
      <template v-if="$slots.actions">
        <v-layout justify-end column>
          <v-divider key="divider" />
          <v-card-actions
            key="actions"
            class="justify-end align-center pa-3"
          >
            <slot name="actions" />
          </v-card-actions>
        </v-layout>
      </template>
    </template>
  </v-card>
</template>

<script>
import { computed, inject } from 'vue';

import { CSS_COLORS_VARS } from '@/config';

import { useModals } from '@/hooks/modals';

import ModalTitleButtons from './modal-title-buttons.vue';
import ModalMassActionsPanel from './modal-mass-actions-panel.vue';

export default {
  inject: ['$modal'],
  components: { ModalTitleButtons, ModalMassActionsPanel },
  props: {
    fillHeight: {
      type: Boolean,
      default: false,
    },
    minimize: {
      type: Boolean,
      default: false,
    },
    close: {
      type: [Boolean, Function],
      default: false,
    },
    titleColor: {
      type: String,
      default: CSS_COLORS_VARS.primary,
    },
    titleClass: {
      type: String,
      required: false,
    },
    textClass: {
      type: String,
      required: false,
    },
  },
  setup(props) {
    const modal = inject('$modal');
    const modals = useModals();

    const isAutoHeight = computed(() => Boolean(modals.dialogPropsMap?.[modal.name]?.autoHeight));

    const titleStyle = computed(() => ({
      backgroundColor: props.titleColor,
    }));

    const wrapperClasses = computed(() => ({
      'fill-min-height': props.fillHeight,
      'modal-wrapper--auto-height': isAutoHeight.value,
    }));

    return {
      titleStyle,
      wrapperClasses,
    };
  },
};
</script>

<style lang="scss" scoped>
.modal-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  overflow: hidden;

  &--auto-height {
    height: auto;
    overflow: visible;

    .v-card__text {
      flex: 0 1 auto;
      min-height: auto;
      overflow-y: visible;
    }
  }

  &__title {
    display: flex;
    justify-content: space-between;
    width: 100%;
    align-items: center;

    & > div {
      display: flex;
    }
  }

  &:not(.modal-wrapper--auto-height) .v-card__text {
    flex: 1 1 0;
    min-height: 0;
    overflow-y: auto;
  }

  & > .layout {
    flex: 0 0 auto;
  }
}
</style>
