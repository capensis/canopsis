<template>
  <c-form-block-row :label="$t('metaAlarmRule.threshold')" class="meta-alarm-rule-threshold-field">
    <c-label
      :label="$t('metaAlarmRule.threshold')"
      :help-text="$t('metaAlarmRule.thresholdHelpText')"
    />
    <v-radio-group v-field="value.threshold_type">
      <v-radio :value="META_ALARMS_THRESHOLD_TYPES.thresholdCount" color="primary">
        <template #label>
          <v-layout class="gap-3" align-center>
            {{ $t('common.count') }}
            <c-number-field
              v-if="value.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdCount"
              v-field="value.threshold_count"
              :min="0"
              name="thresholdCount"
            />
          </v-layout>
        </template>
      </v-radio>

      <v-radio :value="META_ALARMS_THRESHOLD_TYPES.thresholdRate" color="primary">
        <template #label>
          <v-layout class="gap-3" align-center>
            {{ $t('common.rate') }}
            <c-percents-field
              v-if="value.threshold_type === META_ALARMS_THRESHOLD_TYPES.thresholdRate"
              v-field="value.threshold_rate"
              name="thresholdRate"
            />
          </v-layout>
        </template>
      </v-radio>
    </v-radio-group>
  </c-form-block-row>
</template>

<script>
import { META_ALARMS_THRESHOLD_TYPES } from '@/constants';

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
.meta-alarm-rule-threshold-field {
  .v-radio {
    width: 100%;
    height: 40px;
  }

  .v-text-field {
    margin-top: 0;
    padding-top: 0;
  }
}
</style>
