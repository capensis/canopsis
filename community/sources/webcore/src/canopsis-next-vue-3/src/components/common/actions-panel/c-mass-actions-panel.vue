<template>
  <Teleport v-if="mounted && targetElement" :to="targetSelector">
    <div :style="style" class="c-mass-actions-panel">
      <c-layout align="center" justify="center">
        <v-card class="c-mass-actions-panel__card">
          <v-card-text class="pa-1">
            <c-layout class="pa-2" gap="3" align="center">
              <span class="c-mass-actions-panel__message mr-1">
                {{ message }}
              </span>
              <v-divider vertical />
              <c-enabled-field
                :model-value="modelValue"
                :label="t('common.massActionsPanel.keepSelectedAfterAction')"
                class="c-mass-actions-panel__enabled-field"
                density="compact"
                hide-details
                no-margin
                @update:model-value="$emit('update:modelValue', $event)"
              />
              <v-divider vertical />
              <slot name="actions" />
              <v-divider vertical />
              <c-action-btn
                :tooltip="t('common.massActionsPanel.clearSelection')"
                icon="close"
                small
                @click="$emit('clear:selected')"
              />
            </c-layout>
          </v-card-text>
        </v-card>
      </c-layout>
    </div>
  </Teleport>
</template>

<script>
import { computed, inject, onBeforeUnmount, onMounted, ref } from 'vue';

import { TELEPORTS_NAMES } from '@/constants';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    defaultValue: {
      type: Boolean,
      default: false,
    },
    selected: {
      type: Array,
      default: () => [],
    },
    x: {
      type: Number,
      default: 0,
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

  emits: ['update:modelValue', 'clear:selected'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const mounted = ref(false);
    const modal = inject('$modal', null);

    const targetName = computed(
      () => `${TELEPORTS_NAMES.massActionsPanel}${modal?.id ? `-${modal.id}` : ''}`,
    );
    const targetSelector = computed(() => `[data-teleport-target="${targetName.value}"]`);
    const targetElement = computed(() =>
      mounted.value ? document.querySelector(targetSelector.value) : null,
    );

    const message = computed(() =>
      t('common.massActionsPanel.recordsSelected', props.selected.length, {
        count: props.selected.length,
      }),
    );

    const style = computed(() => ({
      gridColumnStart: props.x + 1,
      gridColumnEnd: props.x + 1 + props.w,
      gridRowStart: props.y + 1,
      gridRowEnd: props.y + props.h + 1,
    }));

    onMounted(() => (mounted.value = true));
    onBeforeUnmount(() => emit('update:modelValue', props.defaultValue));

    return {
      t,
      mounted,
      targetElement,
      targetSelector,
      message,
      style,
    };
  },
};
</script>

<style lang="scss">
.c-mass-actions-panel {
  display: block;
  justify-content: center;
  width: 100%;
  padding-top: 20px;
  z-index: 3;
  position: relative;

  &__card {
    max-width: 100%;
    pointer-events: auto;
  }

  &__message {
    white-space: nowrap;
  }

  &__enabled-field {
    margin: 0;

    .v-switch {
      margin: 0;
      padding: 0;
    }

    .v-input__control,
    .v-selection-control {
      min-height: 32px;
    }

    .v-label {
      font-size: 13px;
    }
  }
}
</style>
