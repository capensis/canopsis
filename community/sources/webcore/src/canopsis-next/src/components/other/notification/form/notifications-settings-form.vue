<template lang="pug">
  v-layout(column)
    h4.subheading {{ $t('notificationSettings.instruction.header') }}
    v-checkbox.my-2(
      v-field="form.instruction.rate",
      :disabled="!form.instruction.rate_frequency.value",
      :label="$t('notificationSettings.instruction.rate')",
      color="primary",
      hide-details
    )
    v-layout
      v-flex.subheading.mt-3(xs2) {{ $t('notificationSettings.instruction.rateFrequency') }}
      v-flex(xs3)
        c-duration-field(
          v-field="form.instruction.rate_frequency",
          :disabled="!form.instruction.rate",
          :required="form.instruction.rate",
          :label="$t('notificationSettings.instruction.duration')",
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
