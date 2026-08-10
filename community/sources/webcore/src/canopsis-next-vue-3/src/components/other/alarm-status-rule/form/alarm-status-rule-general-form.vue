<template>
  <c-layout column>
    <c-name-field
      v-model="nameField"
      :max-length="255"
      autofocus
      required
    />

    <c-form-block>
      <c-form-block-row :label="t('common.priority')">
        <c-priority-field v-model="priorityField" :disabled="defaultRule" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.duration')">
        <c-duration-field v-model:duration="durationField" required />
      </c-form-block-row>

      <c-form-block-row v-if="flapping" :label="t('common.frequencyLimit')">
        <c-number-field
          v-model:value="frequencyLimitField"
          :label="t('common.frequencyLimit')"
          :min="1"
          name="freq_limit"
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.description')">
        <c-description-field v-model="descriptionField" required />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
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
    };
  },
};
</script>
