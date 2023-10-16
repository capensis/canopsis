<template>
  <v-list-group>
    <template #activator="">
      <v-list-item>{{ $t('settings.filterOnOpenResolved') }}</v-list-item>
    </template>
    <v-container>
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
    </v-container>
  </v-list-group>
</template>

<script>
import { ALARMS_OPENED_VALUES } from '@/constants';

export default {
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
