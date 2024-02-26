<template>
  <widget-settings-item :title="$t('settings.availability.graphSettings')">
    <c-enabled-field v-field="form.enabled" />
    <v-expand-transition>
      <div v-if="form.enabled">
        <c-quick-date-interval-type-field
          v-field="form.default_time_range"
          :name="defaultTimeRangeFieldName"
          :ranges="intervalRanges"
          :label="$t('settings.defaultTimeRange')"
        />
        <c-availability-show-type-field
          v-field="form.default_show_type"
          :label="$t('settings.availability.defaultAvailabilityDisplay')"
          :name="defaultShowTypeFieldName"
        />
      </div>
    </v-expand-transition>
  </widget-settings-item>
</template>

<script>
import { AVAILABILITY_QUICK_RANGES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'availability',
    },
  },
  computed: {
    defaultTimeRangeFieldName() {
      return `${this.name}.default_time_range`;
    },

    defaultShowTypeFieldName() {
      return `${this.name}.show_type`;
    },

    intervalRanges() {
      return Object.values(AVAILABILITY_QUICK_RANGES);
    },
  },
};
</script>
