<template lang="pug">
  v-switch(
    v-field="value",
    :label="label || $t('common.enabled')",
    :color="color",
    :disabled="disabled",
    :hide-details="hideDetails"
  )
    template(#append="")
      slot(name="append")
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
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
