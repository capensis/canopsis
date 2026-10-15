<template>
  <div class="d-flex align-center ga-4 pt-2">
    <div :style="labelStyles">
      <div class="c-information-block-row__label text-body-2">
        {{ label }}
      </div>
    </div>
    <div class="text-break">
      <slot>{{ value }}</slot>
    </div>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    value: {
      type: [String, Number, Array, Object],
      required: false,
    },
    width: {
      type: [Number, String],
      default: 200,
    },
  },
  computed: {
    labelStyles() {
      return {
        minWidth: `${this.width}px`,
        maxWidth: `${this.width}px`,
      };
    },
  },
};
</script>

<style scoped>
.c-information-block-row__label {
  font-weight: 700 !important;
}
</style>
