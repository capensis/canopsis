<template>
  <widget-settings-item :title="$t('settings.chart.selectMetrics')">
    <c-alarm-metric-presets-field
      v-field="value"
      :with-color="withColor"
      :with-external="withExternal"
      :with-aggregate-function="withAggregateFunction"
      :parameters="parameters"
      :only-group="onlyGroup"
      :only-external="onlyExternal"
      :min="min"
    />
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

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
    withExternal: {
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
    onlyExternal: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
