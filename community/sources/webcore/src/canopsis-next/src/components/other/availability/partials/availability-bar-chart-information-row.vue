<template>
  <v-layout class="availability-bar-chart-information-row gap-4" align-center>
    <div :style="{ backgroundColor: color }" class="availability-bar-chart-information-row__color" />
    <div class="availability-bar-chart-information-row__label">
      {{ label }}
    </div>
    <div class="availability-bar-chart-information-row__value">
      <slot />
    </div>
  </v-layout>
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

<style lang="scss">
.availability-bar-chart-information-row {
  &__color {
    height: 20px;
    width: 20px;
  }

  &__label {
    flex-shrink: 0;
    width: 250px;
  }

  &__value {}
}
</style>
