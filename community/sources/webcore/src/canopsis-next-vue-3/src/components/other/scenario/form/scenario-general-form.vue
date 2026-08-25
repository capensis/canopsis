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
        <c-priority-field v-model="priorityField" />
      </c-form-block-row>

      <c-form-block-row :label="t('common.trigger', 2)">
        <c-triggers-field
          v-model="triggersField"
          with-additional-values
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.delay')">
        <c-duration-field
          v-model:duration="delayField"
          :label="t('common.delay')"
          name="delay"
          clearable
        />
      </c-form-block-row>

      <c-form-block-row :label="t('common.disableDuringPeriods')">
        <c-disable-during-periods-field v-model="disableDuringPeriodsField" />
      </c-form-block-row>
    </c-form-block>
  </c-layout>
</template>

<script>
import { useFormContext } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  setup() {
    const { t } = useI18n();
    const form = useFormContext();
    const {
      nameField,
      priorityField,
      triggersField,
      delayField,
      disableDuringPeriodsField,
    } = form.models({
      nameField: 'name',
      priorityField: 'priority',
      triggersField: 'triggers',
      delayField: 'delay',
      disableDuringPeriodsField: 'disable_during_periods',
    });

    return {
      t,
      nameField,
      priorityField,
      triggersField,
      delayField,
      disableDuringPeriodsField,
    };
  },
};
</script>
