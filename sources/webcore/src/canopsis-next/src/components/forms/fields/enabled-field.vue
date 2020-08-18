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
      required: false,
    },
    label: {
      type: String,
      required: false,
    },
    color: {
      type: String,
      required: false,
    },
    hideDetails: {
      type: Boolean,
      default: false,
    },
  },
};
</script>
