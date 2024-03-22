<template>
  <widget-settings-item :title="label">
    <c-columns-with-template-field
      v-field="columns"
      :type="type"
      :template="template"
      :templates="templates"
      :templates-pending="templatesPending"
      :with-template="withTemplate"
      :with-html="withHtml"
      :with-color-indicator="withColorIndicator"
      :excluded-columns="excludedColumns"
      @update:template="updateTemplate"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'columns',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
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
    excludedColumns: {
      type: Array,
      required: false,
    },
  },
  methods: {
    updateTemplate(...args) {
      this.$emit('update:template', ...args);
    },
  },
};
</script>
