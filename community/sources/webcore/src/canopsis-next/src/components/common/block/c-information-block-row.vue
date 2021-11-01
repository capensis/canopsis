<template lang="pug">
  v-layout.c-information-block-row.pt-2(row)
    v-flex.c-information-block-row--label
      h4.subheading.font-weight-bold {{ label }}
    v-flex
      slot {{ value }}
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number],
      required: false,
    },
  },
};
</script>

<style lang="scss">
.c-information-block-row {
  &--label {
    min-width: 160px;
    max-width: 160px;
  }
}
</style>
