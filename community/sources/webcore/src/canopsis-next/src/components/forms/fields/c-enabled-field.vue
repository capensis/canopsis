<template lang="pug">
  v-switch(
    v-field="value",
    v-on="$listeners",
    :label="label || $t('common.enabled')",
    :color="color",
    :disabled="disabled",
    :readonly="readonly",
    :hide-details="hideDetails"
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
      type: Boolean,
      default: true,
    },
    label: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    readonly: {
      type: Boolean,
      default: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
