<template>
  <widget-settings-item :title="$t('settings.filterOnOpenResolved')">
    <v-layout>
      <v-radio-group
        class="mt-0"
        v-field="value"
        name="opened"
        hide-details
        mandatory
      >
        <v-radio
          v-for="type in types"
          :key="type.value"
          :label="type.label"
          :value="type.value"
          color="primary"
        />
      </v-radio-group>
    </v-layout>
  </widget-settings-item>
</template>

<script>
import { ALARMS_OPENED_VALUES } from '@/constants';

import WidgetSettingsItem from '@/components/sidebars/partials/widget-settings-item.vue';

export default {
  components: { WidgetSettingsItem },
  props: {
    value: {
      type: Boolean,
      default: ALARMS_OPENED_VALUES.opened,
    },
  },
  computed: {
    types() {
      return Object.values(ALARMS_OPENED_VALUES).map(value => ({
        value,
        label: this.$t(`settings.openedTypes.${value}`),
      }));
    },
  },
};
</script>
