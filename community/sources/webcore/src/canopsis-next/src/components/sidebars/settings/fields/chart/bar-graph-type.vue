<template lang="pug">
  widget-settings-item(:title="$t('settings.chart.graphType')")
    v-radio-group.pt-0.mt-0(v-field="value", :name="name", hide-details)
      v-radio(
        :value="false",
        :label="$t('settings.chart.separateBars')",
        color="primary"
      )
      v-radio(
        :value="true",
        :label="$t('settings.chart.stackedBars')",
        color="primary"
      )
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/settings/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Boolean,
      required: true,
    },
    name: {
      type: String,
      default: 'stacked',
    },
  },
};
</script>
