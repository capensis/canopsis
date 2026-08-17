<template>
  <c-layout class="availability-bar-chart-information-row" gap="4" align="center">
    <div :style="{ backgroundColor: color }" class="availability-bar-chart-information-row__color" />
    <div class="availability-bar-chart-information-row__label">
      {{ label }}
    </div>
    <div>
      <slot />
    </div>
  </c-layout>
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      default: 'transparent',
    },
    label: {
      type: String,
      required: true,
    },
  },
};
</script>

<style lang="scss" scoped>
.availability-bar-chart-information-row {
  &__color {
    width: 20px;
    height: 20px;
  }

  &__label {
    width: 250px;
    flex-shrink: 0;
  }
}
</style>
