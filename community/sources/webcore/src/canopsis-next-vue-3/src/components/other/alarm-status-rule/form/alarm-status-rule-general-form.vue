<template>
  <c-layout column>
    <c-name-field
      v-model="nameField"
      :max-length="FORM_MAX_STRING_LENGTH"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('common.priority')">
        <c-priority-field v-model="priorityField" :disabled="defaultRule" :min="MIN_PRIORITY" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.duration')">
        <c-duration-field v-model:duration="durationField" :min="MIN_DURATION_VALUE" required />
      </c-form-block-row>

      <c-form-block-row v-if="flapping" :label="t('common.frequencyLimit')">
        <c-number-field
          v-model:value="frequencyLimitField"
          :label="t('common.frequencyLimit')"
          :min="MIN_FREQUENCY_LIMIT"
          name="freq_limit"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.description')" indented>
        <c-description-field v-model="descriptionField" required />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import {
  FORM_MAX_STRING_LENGTH,
  MIN_DURATION_VALUE,
  MIN_FREQUENCY_LIMIT,
  MIN_PRIORITY,
} from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    flapping: {
      type: Boolean,
      default: false,
    },
    defaultRule: {
      type: Boolean,
      default: false,
    },
  },

  setup() {
    const { t } = useI18n();

    const { nameField, priorityField, durationField, frequencyLimitField, descriptionField } =
      useFormModels({
        nameField: 'name',
        priorityField: 'priority',
        durationField: 'duration',
        frequencyLimitField: 'freq_limit',
        descriptionField: 'description',
      });

    return {
      t,
      nameField,
      priorityField,
      durationField,
      frequencyLimitField,
      descriptionField,
      FORM_MAX_STRING_LENGTH,
      MIN_DURATION_VALUE,
      MIN_FREQUENCY_LIMIT,
      MIN_PRIORITY,
    };
  },
};
</script>
