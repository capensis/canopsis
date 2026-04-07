<template>
  <portal :to="target">
    <div :style="style" class="c-mass-actions-panel">
      <v-layout align-center justify-center>
        <v-card class="c-mass-actions-panel__card">
          <v-card-text class="pa-1">
            <v-layout class="gap-4 pa-3" align-center>
              <span class="c-mass-actions-panel__message mr-1">{{ message }}</span>
              <v-divider vertical />
              <c-enabled-field
                v-field="value"
                :label="$t('common.massActionsPanel.keepSelectedAfterAction')"
                class="ma-0"
                hide-details
              />
              <v-divider vertical />
              <slot name="actions" />
              <v-divider vertical />
              <c-action-btn
                :tooltip="$t('common.massActionsPanel.clearSelection')"
                icon="close"
                small
                @click="clearSelected"
              />
            </v-layout>
          </v-card-text>
        </v-card>
      </v-layout>
    </div>
  </portal>
</template>

<script>
import { computed, inject } from 'vue';

import { useI18n } from '@/hooks/i18n';

import { PORTALS_NAMES } from '@/constants/common';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      default: false,
    },
    selected: {
      type: Array,
      default: () => [],
    },
    x: {
      type: Number,
      default: 1,
    },
    y: {
      type: Number,
      default: 1,
    },
    w: {
      type: Number,
      default: 12,
    },
    h: {
      type: Number,
      default: 1,
    },
  },
  setup(props, { emit }) {
    const { tc } = useI18n();
    const modal = inject('$modal', null);

    const target = computed(() => `${PORTALS_NAMES.massActionsPanel}${modal?.id ? `-${modal.id}` : ''}`);

    const message = computed(() => (
      tc('common.massActionsPanel.recordsSelected', props.selected.length, { count: props.selected.length })
    ));

    const style = computed(() => (
      {
        gridColumnStart: props.x + 1,
        gridColumnEnd: props.x + 1 + props.w,
        gridRowStart: props.y + 1,
        gridRowEnd: props.y + 2,
      }
    ));

    /**
     * Emits 'clear:selected' to notify the parent to clear the mass selection.
     */
    const clearSelected = () => emit('clear:selected');

    /**
     * Emits 'update:keep-selected' to notify the parent to update the keep selected after action.
     *
     * @param {boolean} value
     */
    const updateKeepSelected = value => emit('update:keep-selected', value);

    return {
      target,
      message,
      style,

      clearSelected,
      updateKeepSelected,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-mass-actions-panel {
  display: block;
  justify-content: center;
  width: 100%;
  padding-top: 20px;
  z-index: 1;
  position: relative;

  &__card {
    max-width: 100%;
    pointer-events: auto;
  }

  &__message {
    white-space: nowrap;
  }
}
</style>
