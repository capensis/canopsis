<template>
  <widget-settings-item :title="$t('settings.chart.fontSize')">
    <c-radio-group
      :model-value="enabled"
      hide-details
      standalone
      @update:model-value="setEnabled"
    >
      <v-radio
        v-for="item in items"
        :key="String(item.value)"
        :value="item.value"
        :label="item.text"
        color="primary"
      />
    </c-radio-group>
    <c-number-field
      v-if="enabled"
      :value="modelValue"
      name="font_size"
      required
      @update:value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import { NUMBERS_CHART_DEFAULT_FONT_SIZE } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: { modelValue: { type: Number, default: undefined } },
  emits: ['update:modelValue'],
  computed: {
    enabled() {
      return this.modelValue !== undefined;
    },
    items() {
      return [
        { value: false, text: this.$t('settings.chart.auto') },
        { value: true, text: this.$t('settings.chart.manual') },
      ];
    },
  },
  methods: {
    setEnabled(value) {
      this.$emit('update:modelValue', value ? NUMBERS_CHART_DEFAULT_FONT_SIZE : undefined);
    },
  },
};
</script>
