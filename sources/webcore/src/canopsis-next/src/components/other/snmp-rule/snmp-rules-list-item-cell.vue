<template lang="pug" functional>
  div
    div(v-for="field in props.fields", :key="field")
      v-chip.va-middle(color="info", text-color="white")
        strong {{ field }}
      v-icon.va-middle(large) arrow_right_alt
      v-chip.va-middle(v-show="props.source[field]", color="warning", text-color="white")
        strong {{ props.source[field] }}
</template>

<script>
export default {
  props: {
    fields: {
      type: Array,
      required: true,
    },
    source: {
      type: Object,
      default: () => ({}),
    },
  },
};
</script>

<style lang="scss" scoped>
  .va-middle {
    vertical-align: middle;
  }
</style>
