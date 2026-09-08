<template>
  <c-form-block>
    <c-form-block-row
      :label="$t('notificationSetting.instruction.rate')"
      align-center
    >
      <c-enabled-field
        v-field="form.instruction.rate"
        :label="$t('notificationSetting.instruction.rate')"
        :disabled="disabled || !form.instruction.rate_frequency.value"
        hide-details
        no-margin
      />
    </c-form-block-row>

    <c-form-block-row :label="$t('notificationSetting.instruction.rateFrequency')">
      <c-duration-field
        v-field="form.instruction.rate_frequency"
        :label="$t('notificationSetting.instruction.duration')"
        :units-label="$t('common.unit')"
        :disabled="disabled || !form.instruction.rate"
        :required="form.instruction.rate"
        :min="1"
        name="instruction.rate_frequency"
      />
    </c-form-block-row>
  </c-form-block>
</template>

<script>
export default {
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
