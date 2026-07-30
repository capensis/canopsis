<template>
  <c-form-block>
    <c-form-block-row :label="t('notificationSetting.instruction.rate')" align-center>
      <c-enabled-field
        v-model="instructionRateField"
        :label="t('notificationSetting.instruction.rate')"
        :disabled="disabled || !instructionRateFrequencyField.value"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row :label="t('notificationSetting.instruction.rateFrequency')">
      <c-duration-field
        v-model:duration="instructionRateFrequencyField"
        :label="t('notificationSetting.instruction.duration')"
        :units-label="t('common.unit')"
        :disabled="disabled || !instructionRateField"
        :required="instructionRateField"
        :min="1"
        name="instruction.rate_frequency"
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  setup() {
    const { t } = useI18n();
    const { instructionRateField, instructionRateFrequencyField } = useFormModels({
      instructionRateField: 'instruction.rate',
      instructionRateFrequencyField: 'instruction.rate_frequency',
    });

    return {
      t,
      instructionRateField,
      instructionRateFrequencyField,
    };
  },
};
</script>
