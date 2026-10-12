<template>
  <widget-settings-item :title="$t('settings.chart.graphType')">
    <c-radio-group
      :model-value="modelValue"
      hide-details
      @update:model-value="$emit('update:modelValue', $event)"
    >
      <v-radio
        v-for="item in items"
        :key="String(item.value)"
        :value="item.value"
        :label="item.text"
        color="primary"
      />
    </c-radio-group>
  </widget-settings-item>
</template>

<script>
import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: { modelValue: { type: Boolean, default: false } },
  emits: ['update:modelValue'],
  computed: {
    items() {
      return [
        { value: false, text: this.$t('settings.chart.separateBars') },
        { value: true, text: this.$t('settings.chart.stackedBars') },
      ];
    },
  },
};
</script>
