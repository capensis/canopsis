<template>
  <div :class="wrapperClass" class="c-enabled-field">
    <v-switch
      :model-value="currentValue"
      :class="{ 'ma-0': withBackground || noMargin, 'pa-0': noMargin }"
      :label="label || t('common.enabled')"
      :color="color"
      :disabled="disabled"
      :readonly="readonly"
      :hide-details="hideDetails || withBackground"
      v-bind="$attrs"
      @update:model-value="updateModel"
    >
      <template v-if="$slots.label" #label>
        <slot name="label" />
      </template>

      <template v-if="$slots.append" #append>
        <slot name="append" />
      </template>
    </v-switch>
  </div>
</template>

<script>
import { computed } from 'vue';

import { useI18n } from '@/hooks/i18n';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: Boolean,
      default: true,
    },
    value: {
      type: Boolean,
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
    withBackground: {
      type: Boolean,
      default: false,
    },
    noMargin: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const currentValue = computed(() => props.modelValue ?? props.value);

    const wrapperClass = computed(() =>
      props.withBackground
        ? {
            'c-enabled-field--with-background': true,
            'c-enabled-field--with-background--disabled': !currentValue.value,
          }
        : {},
    );

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    return {
      t,
      currentValue,
      wrapperClass,
      updateModel,
    };
  },
};
</script>

<style lang="scss" scoped>
.c-enabled-field {
  display: flex;
}

.c-enabled-field--with-background {
  padding: 0 8px;
  border-radius: 5px;
  border: 1px solid rgb(var(--v-theme-success));
  background-color: rgba(var(--v-theme-success), 0.12);

  :deep(.v-input) {
    margin: 0;
    padding: 0;
  }

  &--disabled {
    border-color: rgb(var(--v-theme-error));
    background-color: rgba(var(--v-theme-error), 0.12);
  }
}
</style>
