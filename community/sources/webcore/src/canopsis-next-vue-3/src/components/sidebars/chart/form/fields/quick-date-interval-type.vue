<template>
  <widget-settings-item :title="$t('settings.defaultTimeRange')">
    <c-quick-date-interval-type-field
      :value="modelValue"
      :name="name"
      :ranges="intervalRanges"
      @update:value="$emit('update:modelValue', $event)"
    />
  </widget-settings-item>
</template>

<script>
import { QUICK_RANGES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    modelValue: {
      type: String,
      default: undefined,
    },
    name: {
      type: String,
      default: '',
    },
    ranges: {
      type: Object,
      default: () => QUICK_RANGES,
    },
  },
  emits: ['update:modelValue'],
  computed: {
    intervalRanges() {
      const { [QUICK_RANGES.custom.value]: custom, ...ranges } = this.ranges;

      return Object.values(ranges);
    },
  },
};
</script>
