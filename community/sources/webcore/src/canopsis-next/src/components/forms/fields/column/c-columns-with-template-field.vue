<template>
  <c-layout column>
    <c-widget-template-field
      :model-value="template"
      :templates="templates"
      :pending="templatesPending"
      @update:model-value="updateTemplate"
    />
    <span class="text-body-2 my-2">
      {{ $t('common.column', 2) }}
    </span>
    <c-columns-field v-bind="$attrs" :model-value="modelValue" @update:model-value="updateValue" />
  </c-layout>
</template>

<script>
import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

export default {
  inheritAttrs: false,

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    template: {
      type: [String, Symbol],
      required: false,
    },
    templates: {
      type: Array,
      default: () => [],
    },
    templatesPending: {
      type: Boolean,
      default: false,
    },
  },

  emits: ['update:modelValue', 'update:template'],

  setup(props, { emit }) {
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'columns', emit);

    return {
      updateTemplate,
      updateValue,
    };
  },
};
</script>
