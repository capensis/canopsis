<template>
  <c-textarea
    v-bind="$attrs"
    :model-value="modelValue"
    :label="label || t('common.description')"
    :name="name"
    :disabled="disabled"
    :required="required"
    :max-length="maxLength"
    @update:model-value="updateModel"
  >
    <template v-if="helpText" #append>
      <c-help-icon :text="helpText" icon="help" left />
    </template>
  </c-textarea>
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
      default: 'description',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    helpText: {
      type: String,
      default: '',
    },
    maxLength: {
      type: Number,
      required: false,
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
