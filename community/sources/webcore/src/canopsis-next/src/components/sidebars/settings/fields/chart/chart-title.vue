<template lang="pug">
  widget-settings-item(:title="$tc('common.header')", optional)
    v-text-field.mt-0.pt-0(
      v-field="value",
      :placeholder="$t('settings.headerTitle')",
      :name="name"
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      default: 'parameters.chart_title',
    },
  },
};
</script>
