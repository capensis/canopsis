<template>
  <c-layout class="pt-2" align="center">
    <div :style="labelStyles">
      <div class="text-subtitle-1 font-weight-bold">
        {{ label }}
      </div>
    </div>
    <div class="text-break">
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
      type: [String, Number, Array, Object],
      required: false,
    },
    width: {
      type: [Number, String],
      default: 160,
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
