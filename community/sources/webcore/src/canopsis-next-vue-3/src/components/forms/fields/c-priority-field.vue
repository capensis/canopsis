<template>
  <c-number-field
    v-bind="$attrs"
    :value="modelValue"
    :label="label || t('common.priority')"
    :name="name"
    :disabled="disabled"
    @update:value="updateModel"
  />
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: [Number, String],
      default: undefined,
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'priority',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const updateModel = (value) => {
      emit('update:modelValue', value);
    };

    return {
      t,
      updateModel,
    };
  },
};
</script>
