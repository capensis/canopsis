<template>
  <c-select
    :model-value="selected"
    :items="templatesWithCustom"
    :label="t('common.template', 1)"
    :loading="pending"
    :clearable="clearable"
    item-title="text"
    item-value="value"
    return-object
    @update:model-value="selectTemplate"
  />
</template>

<script>
import { computed } from 'vue';

import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    modelValue: {
      type: [String, Object],
      required: false,
    },
    value: {
      type: [String, Object],
      required: false,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    pending: {
      type: Boolean,
      default: false,
    },
    valueKey: {
      type: String,
      default: 'value',
    },
    clearable: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue'],

  setup(props, { emit }) {
    const { t } = useI18n();

    const currentValue = computed(() => props.modelValue ?? props.value);
    const templatesWithCustom = computed(() => [
      {
        value: CUSTOM_WIDGET_TEMPLATE,
        text: t('common.custom'),
        content: '',
      },

      ...props.templates.map(item => ({
        ...item,
        value: item._id,
        text: item.title,
        content: item.content ?? item.text ?? '',
      })),
    ]);

    const selected = computed(
      () =>
        templatesWithCustom.value.find(template => template.value === currentValue.value) ??
        templatesWithCustom.value[0],
    );

    const selectTemplate = (template) => {
      emit('update:modelValue', template);
    };

    return {
      t,
      templatesWithCustom,
      selected,
      selectTemplate,
    };
  },
};
</script>
