<template>
  <widget-settings-item :title="$t('settings.defaultTimeRange')">
    <c-quick-date-interval-type-field
      v-field="value"
      :name="name"
      :ranges="intervalRanges"
    />
  </widget-settings-item>
</template>

<script>
import { QUICK_RANGES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: String,
      required: false,
    },
    name: {
      type: String,
      required: false,
    },
    ranges: {
      type: Object,
      default: () => QUICK_RANGES,
    },
  },
  computed: {
    intervalRanges() {
      const {
        [QUICK_RANGES.custom.value]: custom,
        ...ranges
      } = this.ranges;

      return Object.values(ranges);
    },
  },
};
</script>
