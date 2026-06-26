<template>
  <c-text-field
    v-bind="$attrs"
    :model-value="modelValue"
    :label="label || t('common.name')"
    :name="name"
    :required="required"
    :max-length="maxLength"
    :regex="regex"
    @update:model-value="updateModel"
  >
    <template
      v-if="tooltip"
      #append
    >
      <c-help-icon
        :text="tooltip"
        icon="help"
        left
      />
    </template>

    <template #append-outer>
      <slot name="append-outer" />
    </template>
  </c-text-field>
</template>

<script>
import { useI18n } from '@/hooks/i18n';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'name',
    },
    required: {
      type: Boolean,
      default: false,
    },
    maxLength: {
      type: Number,
      required: false,
    },
    regex: {
      type: RegExp,
      required: false,
    },
    tooltip: {
      type: String,
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const updateModel = value => {
      emit('update:modelValue', value);
    };

    return {
      t,
      updateModel,
    };
  },
};
</script>
