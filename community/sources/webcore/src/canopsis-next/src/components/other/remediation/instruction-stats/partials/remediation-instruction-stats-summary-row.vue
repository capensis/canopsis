<template>
  <c-layout class="mb-3 remediation-instruction-stats-summary-row">
    <div class="remediation-instruction-stats-summary-row--label">
      <span class="font-weight-bold">
        {{ label }}
      </span>
    </div>
    <div>
      <slot>{{ value }}</slot>
    </div>
  </c-layout>
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
.remediation-instruction-stats-summary-row {
  &--label {
    max-width: 200px;
    width: 200px;
  }
}
</style>
