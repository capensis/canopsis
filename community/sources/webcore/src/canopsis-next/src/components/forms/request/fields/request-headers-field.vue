<template lang="pug">
  c-text-pairs-field(
    v-field="value",
    :text-label="$t('common.headerKey')",
    :value-label="$t('common.headerValue')",
    :name="name",
    :disabled="disabled",
    :hints="headersHints",
    text-required
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
      type: Array,
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
  computed: {
    headersHints() {
      return [
        {
          text: 'Content-Type',
          value: ['application/json', 'application/javascript', 'application/x-www-form-urlencoded'],
        },
      ];
    },
  },
};
</script>
