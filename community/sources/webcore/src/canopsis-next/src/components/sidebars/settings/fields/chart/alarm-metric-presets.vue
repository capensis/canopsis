<template lang="pug">
  widget-settings-item(:title="$t('settings.chart.selectMetrics')")
    c-alarm-metric-presets-field(
      v-field="value",
      :with-color="withColor",
      :with-aggregate-function="withAggregateFunction",
      :parameters="parameters",
      :only-group="onlyGroup",
      :min="min"
    )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Array,
      required: true,
    },
    withColor: {
      type: Boolean,
      default: false,
    },
    withAggregateFunction: {
      type: Boolean,
      default: false,
    },
    parameters: {
      type: Array,
      required: false,
    },
    onlyGroup: {
      type: Boolean,
      default: false,
    },
    min: {
      type: Number,
      required: false,
    },
  },
};
</script>
