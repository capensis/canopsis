<template lang="pug" functional>
  div
    v-layout(v-for="field in props.fields", :key="field", row)
      v-chip.chip-small.va-middle(color="info", text-color="white")
        span.caption.font-weight-medium {{ field }}
      v-icon.va-middle arrow_right_alt
      v-chip.chip-small.va-middle(v-show="props.source[field]", color="warning", text-color="white")
        span.caption.font-weight-medium {{ props.source[field] }}
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

  .chip-small {
    height: 22px;

    & /deep/ .v-chip__content {
      padding: 0 7px;
    }
  }
</style>
