<template lang="pug">
  v-layout(column)
    h4.subheading {{ $t('notificationSetting.instruction.header') }}
    v-checkbox.my-2(
      v-field="form.instruction.rate",
      :disabled="!form.instruction.rate_frequency.value",
      :label="$t('notificationSetting.instruction.rate')",
      color="primary",
      hide-details
    )
    v-layout
      v-flex.subheading.mt-3(xs2) {{ $t('notificationSetting.instruction.rateFrequency') }}
      v-flex(xs3)
        c-duration-field(
          v-field="form.instruction.rate_frequency",
          :disabled="!form.instruction.rate",
          :required="form.instruction.rate",
          :label="$t('notificationSetting.instruction.duration')",
          long
        )
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
      required: true,
    },
  },
};
</script>
