<template lang="pug">
  div
    template(v-if="source")
      v-layout(v-for="field in fields", :key="field", row)
        v-chip.chip-small.va-middle(color="info", text-color="white")
          span {{ field }}
        v-icon.va-middle arrow_right_alt
        v-chip.chip-small.va-middle(v-show="source[field]", color="warning", text-color="white")
          span {{ source[field] }}
    template(v-else)
      v-layout(v-for="field in fields", :key="field", row)
        v-chip.chip-small.va-middle(color="info", text-color="white")
          span {{ field }}
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
      default: null,
    },
  },
};
</script>

<style lang="scss" scoped>
  .va-middle {
    vertical-align: middle;
  }

  .chip-small {
    height: 20px;

    & ::v-deep .v-chip__content {
      padding: 0 5px;
    }

    span {
      font-size: .75em;
      font-weight: 500;
    }
  }
</style>
