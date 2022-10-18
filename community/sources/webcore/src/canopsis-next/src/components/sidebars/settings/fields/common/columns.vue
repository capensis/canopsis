<template lang="pug">
  widget-settings-item(:title="label")
    c-columns-field(
      v-field="columns",
      :with-template="withTemplate",
      :with-html="withHtml",
      :with-color-indicator="withColorIndicator"
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

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
    columns: {
      type: [Array, Object],
      default: () => [],
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
};
</script>
