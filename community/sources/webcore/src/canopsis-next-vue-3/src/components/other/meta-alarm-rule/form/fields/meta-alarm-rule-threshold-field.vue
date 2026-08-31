<template>
  <c-form-block-row :label="t('metaAlarmRule.threshold')" class="meta-alarm-rule-threshold-field">
    <c-layout gap="2" column>
      <c-label
        :label="t('metaAlarmRule.threshold')"
        :help-text="t('metaAlarmRule.thresholdHelpText')"
      />

      <c-radio-group
        :model-value="config.threshold_type"
        name="threshold_type"
        hide-details
        @update:model-value="updateField('threshold_type', $event)"
      >
        <v-radio
          v-for="option in thresholdOptions"
          :key="option.value"
          :value="option.value"
        >
          <template #label>
            <c-layout class="meta-alarm-rule-threshold-field__option" gap="3" align="center">
              <span>{{ option.text }}</span>

              <c-number-field
                v-if="option.value === config.threshold_type"
                :value="getThresholdValue(option.value)"
                :min="0"
                :max="option.value === META_ALARMS_THRESHOLD_TYPES.thresholdRate ? 100 : undefined"
                :name="getThresholdName(option.value)"
                hide-details
                @update:value="updateThreshold(option.value, $event)"
              />
            </c-layout>
          </template>
        </v-radio>
      </c-radio-group>
    </c-layout>
  </c-form-block-row>
</template>

<script>
import { computed } from 'vue';

import { META_ALARMS_THRESHOLD_TYPES } from '@/constants';

import { useModelField } from '@/hooks/form/model-field';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    config: {
      type: Object,
      default: () => ({}),
    },
  },

  emits: ['update:config'],

  setup(props, { emit }) {
    const { t } = useI18n();
    const { updateField } = useModelField(props, emit);

    const thresholdOptions = computed(() => [
      {
        value: META_ALARMS_THRESHOLD_TYPES.thresholdCount,
        text: t('common.count'),
      },
      {
        value: META_ALARMS_THRESHOLD_TYPES.thresholdRate,
        text: t('common.rate'),
      },
    ]);

    const getThresholdName = type => type === META_ALARMS_THRESHOLD_TYPES.thresholdRate ? 'threshold_rate' : 'threshold_count';

    const getThresholdValue = type => props.config[getThresholdName(type)];
    const updateThreshold = (type, value) => updateField(getThresholdName(type), value);

    return {
      META_ALARMS_THRESHOLD_TYPES,
      t,
      thresholdOptions,
      updateField,
      getThresholdName,
      getThresholdValue,
      updateThreshold,
    };
  },
};
</script>

<style lang="scss" scoped>
.meta-alarm-rule-threshold-field__option {
  width: 100%;

  :deep(.c-number-field) {
    max-width: 120px;
  }
}
</style>
