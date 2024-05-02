<template>
  <v-layout class="meta-alarm-rule-threshold-form">
    <v-radio-group v-field="threshold.threshold_type">
      <v-radio :value="META_ALARMS_THRESHOLD_TYPES.thresholdCount" color="primary">
        <template #label>
          <v-layout class="gap-3" align-center>
            {{ $t('metaAlarmRule.thresholdCount') }}
            <c-number-field
              v-if="threshold.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount"
              v-field="threshold.threshold_count"
              :min="0"
              name="thresholdCount"
            />
            <c-help-icon
              :text="$t('metaAlarmRule.thresholdCountHelpText')"
              icon="help"
              max-width="300"
              top
            />
          </v-layout>
        </template>
      </v-radio>
      <v-radio :value="META_ALARMS_THRESHOLD_TYPES.thresholdRate" color="primary">
        <template #label>
          <v-layout class="gap-3" align-center>
            {{ $t('metaAlarmRule.thresholdRate') }}
            <c-percents-field
              v-if="threshold.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdRate"
              v-field="threshold.threshold_rate"
              name="thresholdRate"
            />
            <c-help-icon
              :text="$t('metaAlarmRule.thresholdRateHelpText')"
              icon="help"
              max-width="300"
              top
            />
          </v-layout>
        </template>
      </v-radio>
    </v-radio-group>
  </v-layout>
</template>

<script>
import { META_ALARMS_THRESHOLD_TYPES } from '@/constants';

export default {
  model: {
    prop: 'threshold',
    event: 'input',
  },
  props: {
    threshold: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    return {
      META_ALARMS_THRESHOLD_TYPES,
    };
  },
};
</script>

<style lang="scss">
.meta-alarm-rule-threshold-form {
  .v-radio {
    width: 100%;
    height: 70px;
  }
}
</style>
