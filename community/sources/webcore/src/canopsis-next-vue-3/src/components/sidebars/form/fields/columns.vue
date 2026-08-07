<template>
  <widget-settings-item
    :title="label"
    :eager="false"
  >
    <c-columns-with-template-field
      :model-value="modelValue"
      :type="type"
      :template="template"
      :templates="templates"
      :templates-pending="templatesPending"
      :with-template="withTemplate"
      :with-html="withHtml"
      :with-color-indicator="withColorIndicator"
      :with-filter-on-click="withFilterOnClick"
      :items="items"
      :excluded-columns="excludedColumns"
      @update:model-value="$emit('update:modelValue', $event)"
      @update:template="updateTemplate"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },

  props: {
    modelValue: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      required: true,
    },
    type: {
      type: String,
      required: false,
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
    withFilterOnClick: {
      type: Boolean,
      default: false,
    },
    items: {
      type: Array,
      required: false,
    },
    excludedColumns: {
      type: Array,
      required: false,
    },
  },

  emits: ['update:modelValue', 'update:template'],

  methods: {
    updateTemplate(...args) {
      this.$emit('update:template', ...args);
    },
  },
};
</script>
