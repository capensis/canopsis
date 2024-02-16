<template>
  <v-layout
    class="py-4"
    column
  >
    <state-setting-information-row
      :label="$t('common.name')"
      :value="stateSetting.type"
    />
    <state-setting-information-row
      :label="$t('common.method')"
      :value="$t(`stateSetting.methods.${stateSetting.method}`)"
    />
    <v-layout
      class="mt-2"
      v-if="isWorstOfShareMethod"
    >
      <v-flex xs4>
        <state-setting-threshold
          :label="$t('common.skipped')"
          :threshold="stateSetting.junit_thresholds.skipped"
        />
      </v-flex>
      <v-flex xs4>
        <state-setting-threshold
          :label="$tc('common.error', 2)"
          :threshold="stateSetting.junit_thresholds.errors"
        />
      </v-flex>
      <v-flex xs4>
        <state-setting-threshold
          :label="$t('common.failures')"
          :threshold="stateSetting.junit_thresholds.failures"
        />
      </v-flex>
    </v-layout>
  </v-layout>
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

import StateSettingInformationRow from './state-setting-information-row.vue';
import StateSettingThreshold from './state-setting-threshold.vue';

export default {
  components: { StateSettingInformationRow, StateSettingThreshold },
  props: {
    stateSetting: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isWorstOfShareMethod() {
      return this.stateSetting.method === STATE_SETTING_METHODS.worstOfShare;
    },
  },
};
</script>
