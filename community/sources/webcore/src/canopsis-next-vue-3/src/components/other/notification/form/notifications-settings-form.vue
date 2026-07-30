<template>
  <c-layout column>
    <h4 class="text-subtitle-1">
      {{ t('notificationSetting.instruction.header') }}
    </h4>

    <c-checkbox
      v-model="instructionRateField"
      :disabled="!instructionRateFrequencyField.value"
      :label="t('notificationSetting.instruction.rate')"
      class="my-2"
      color="primary"
      hide-details
    />

    <c-layout align="start" gap="3">
      <div class="text-subtitle-1 mt-3 notifications-settings-form__label">
        {{ t('notificationSetting.instruction.rateFrequency') }}
      </div>

      <div class="notifications-settings-form__duration">
        <c-duration-field
          v-model:duration="instructionRateFrequencyField"
          :disabled="!instructionRateField"
          :required="instructionRateField"
          :label="t('notificationSetting.instruction.duration')"
          name="instruction.rate_frequency"
          long
        />
      </div>
    </c-layout>
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
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

<style lang="scss" scoped>
.notifications-settings-form__label {
  flex: 0 0 16.6667%;
  max-width: 16.6667%;
}

.notifications-settings-form__duration {
  flex: 0 0 25%;
  max-width: 25%;
}
</style>
