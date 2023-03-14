<template lang="pug">
  v-layout(column)
    v-select(
      :value="template",
      :items="templatesWithCustom",
      :label="$t('common.template')",
      :loading="templatesPending",
      return-object,
      @input="updateTemplate"
    )
    span.body-2.my-2 {{ $tc('common.column', 2) }}
    c-columns-field(
      :columns="columns",
      :with-template="withTemplate",
      :with-html="withHtml",
      :with-color-indicator="withColorIndicator",
      :type="type",
      @input="updateColumns"
    )
</template>

<script>
import { CUSTOM_WIDGET_TEMPLATE } from '@/constants';

import { formBaseMixin } from '@/mixins/form';

import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  mixins: [formBaseMixin],
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    type: {
      type: String,
      required: true,
    },
    columns: {
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
    withTemplate: {
      type: Boolean,
      default: false,
    },
    withHtml: {
      type: Boolean,
      default: false,
    },
    withColorIndicator: {
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
