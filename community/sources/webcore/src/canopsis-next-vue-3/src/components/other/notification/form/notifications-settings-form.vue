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

    <v-row class="flex-grow-0" align="start" density="comfortable">
      <v-col class="text-subtitle-1 mt-3" cols="2">
        {{ t('notificationSetting.instruction.rateFrequency') }}
      </v-col>

      <v-col cols="3">
        <c-duration-field
          v-model:duration="instructionRateFrequencyField"
          :disabled="!instructionRateField"
          :required="instructionRateField"
          :label="t('notificationSetting.instruction.duration')"
          name="instruction.rate_frequency"
          long
        />
      </v-col>
    </v-row>
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
