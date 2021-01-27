<template lang="pug">
  v-switch(
    v-field="value",
    :label="label || $t('common.enabled')",
    :color="color",
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
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
