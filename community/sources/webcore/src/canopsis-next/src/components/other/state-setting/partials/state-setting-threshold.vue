<template>
  <v-layout column>
    <state-setting-information-row
      :label="$t('stateSetting.criterion')"
      :value="label"
    />
    <state-setting-information-row
      :label="$t('stateSetting.states.minor')"
      :value="threshold.minor"
      :unit="unit"
    />
    <state-setting-information-row
      :label="$t('stateSetting.states.major')"
      :value="threshold.major"
      :unit="unit"
    />
    <state-setting-information-row
      :label="$t('stateSetting.states.critical')"
      :value="threshold.critical"
      :unit="unit"
    />
  </v-layout>
</template>

<script>
import { STATE_SETTING_THRESHOLD_TYPES } from '@/constants';

import StateSettingInformationRow from './state-setting-information-row.vue';

export default {
  components: { StateSettingInformationRow },
  props: {
    threshold: {
      type: Object,
      default: () => ({}),
    },
    label: {
      type: String,
      required: true,
    },
  },
  computed: {
    unit() {
      return this.threshold.type === STATE_SETTING_THRESHOLD_TYPES.percent ? '%' : '';
    },
  },
};
</script>
