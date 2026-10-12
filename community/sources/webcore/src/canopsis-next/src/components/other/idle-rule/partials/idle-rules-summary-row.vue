<template>
  <c-layout class="mb-3 c-idle-rules-summary-row">
    <div class="c-idle-rules-summary-row__label">
      <span class="font-weight-bold">
        {{ label }}
      </span>
    </div>
    <div class="c-idle-rules-summary-row__value">
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

<style lang="scss" scoped>
.c-idle-rules-summary-row__label {
  flex: 0 0 200px;
  max-width: 200px;
}

.c-idle-rules-summary-row__value {
  flex: 1 1 auto;
  min-width: 0;
}
</style>
