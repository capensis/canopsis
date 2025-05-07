<template>
  <v-layout column>
    <c-widget-template-field
      :value="template"
      :templates="templates"
      :pending="templatesPending"
      @input="updateTemplate"
    />
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <c-columns-field
      v-bind="$attrs"
      @input="updateValue"
    />
  </v-layout>
</template>

<script>
import { useWidgetTemplateField } from '@/hooks/widget/widget-template';

export default {
  inheritAttrs: false,
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
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
  setup(props, { emit }) {
    const { updateTemplate, updateValue } = useWidgetTemplateField(props, 'columns', emit);

    return {
      updateTemplate,
      updateValue,
    };
  },
};
</script>
