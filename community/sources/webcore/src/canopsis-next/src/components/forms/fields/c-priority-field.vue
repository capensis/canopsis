<template>
  <c-number-field
    v-bind="$attrs"
    :value="modelValue"
    :label="label || t('common.priority')"
    :name="name"
    :min="min"
    :required="required"
    :disabled="disabled"
    @update:value="updateModel"
  >
    <template v-for="(_, slot) in $slots" #[slot]="slotProps">
      <slot :name="slot" v-bind="slotProps || {}" />
    </template>
  </c-number-field>
</template>

<script>
import { MIN_PRIORITY } from '@/constants';

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
    min: {
      type: Number,
      default: MIN_PRIORITY,
    },
    required: {
      type: Boolean,
      default: false,
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
