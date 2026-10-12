<template>
  <v-layout column>
    <v-layout>
      <v-flex xs4>
        <junit-state-setting-method-field v-field="form.method" />
      </v-flex>
    </v-layout>
    <v-layout
      v-if="isWorstOfShareMethod"
      column
    >
      <v-layout>
        <h4 class="subheading font-weight-bold">
          {{ $t('stateSetting.junit.worstLabel') }}
        </h4>
        <c-help-icon
          :text="$t('stateSetting.junit.worstHelpText')"
          icon-class="ml-2"
          max-width="220"
          right
        />
      </v-layout>
      <junit-state-setting-threshold-field
        v-field="form.junit_thresholds.skipped"
        :label="$t('common.skipped')"
        class="pl-4 pt-2"
        name="junit_thresholds.skipped"
      />
      <junit-state-setting-threshold-field
        v-field="form.junit_thresholds.errors"
        :label="$tc('common.error', 2)"
        class="pl-4 pt-2"
        name="junit_thresholds.errors"
      />
      <junit-state-setting-threshold-field
        v-field="form.junit_thresholds.failures"
        :label="$t('common.failures')"
        class="pl-4 pt-2"
        name="junit_thresholds.failures"
      />
    </v-layout>
  </v-layout>
</template>

<script>
import { JUNIT_STATE_SETTING_METHODS } from '@/constants';

import JunitStateSettingMethodField from './fields/junit-state-setting-method-field.vue';
import JunitStateSettingThresholdField from './fields/junit-state-setting-thresholds-field.vue';

export default {
  components: {
    JunitStateSettingMethodField,
    JunitStateSettingThresholdField,
  },
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  computed: {
    isWorstOfShareMethod() {
      return this.form.method === JUNIT_STATE_SETTING_METHODS.worstOfShare;
    },
  },
};
</script>
