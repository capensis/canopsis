<template lang="pug">
  v-layout(column)
    v-layout(row)
      v-flex(xs4)
        state-setting-method-field(v-field="form.method")
    v-layout(v-if="isWorstOfShareMethod", column)
      v-layout(row)
        h4.subheading.font-weight-bold {{ $t('stateSetting.worstLabel') }}
        c-help-icon(:text="$t('stateSetting.worstHelpText')", icon-class="ml-2", max-width="220", right)
      state-setting-threshold-field.pl-4.pt-2(
        v-field="form.junit_thresholds.skipped",
        :label="$t('common.skipped')",
        name="junit_thresholds.skipped"
      )
      state-setting-threshold-field.pl-4.pt-2(
        v-field="form.junit_thresholds.errors",
        :label="$tc('common.error', 2)",
        name="junit_thresholds.errors"
      )
      state-setting-threshold-field.pl-4.pt-2(
        v-field="form.junit_thresholds.failures",
        :label="$t('common.failures')",
        name="junit_thresholds.failures"
      )
</template>

<script>
import { STATE_SETTING_METHODS } from '@/constants';

import { formMixin } from '@/mixins/form';

import StateSettingMethodField from './fields/state-setting-method-field.vue';
import StateSettingThresholdField from './fields/state-setting-thresholds-field.vue';

export default {
  inject: ['$validator'],
  components: { StateSettingThresholdField, StateSettingMethodField },
  mixins: [formMixin],
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
      return this.form.method === STATE_SETTING_METHODS.worstOfShare;
    },
  },
};
</script>
