<template lang="pug">
  v-layout(row)
    v-flex(xs4)
      c-enabled-field(
        v-field="value.enabled",
        :label="$t('scenario.withAuth')",
        :disabled="disabled"
      )
    template(v-if="value.enabled")
      v-flex.mr-3(xs4)
        c-name-field(
          v-field="value.username",
          :label="$t('common.username')",
          :name="`${name}.username`",
          :disabled="disabled"
        )
      v-flex(xs4)
        c-password-field(
          v-field="value.password",
          :name="`${name}.password`",
          :disabled="disabled"
        )
</template>

<script>
export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'auth',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
