<template>
  <v-layout column>
    <v-select
      :value="template"
      :items="templatesWithCustom"
      :label="$t('common.template')"
      :loading="templatesPending"
      return-object
      @input="updateTemplate"
    />
    <span class="text-body-2 my-2">{{ $tc('common.column', 2) }}</span>
    <c-columns-field
      v-bind="$attrs"
      @input="updateColumns"
    />
  </v-layout>
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

export default {
  mixins: [formBaseMixin],
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
  computed: {
    templatesWithCustom() {
      return [
        { value: CUSTOM_WIDGET_TEMPLATE, text: this.$t('common.custom'), columns: [] },

        ...this.templates.map(template => ({
          ...template,

          value: template._id,
          text: template.title,
        })),
      ];
    },
  },
  methods: {
    updateColumns(columns) {
      if (this.template !== CUSTOM_WIDGET_TEMPLATE) {
        this.$emit('update:template', CUSTOM_WIDGET_TEMPLATE, columns);

        return;
      }

      this.updateModel(columns);
    },

    updateTemplate({ value, columns }) {
      this.$emit('update:template', value, columns);
    },
  },
};
</script>
