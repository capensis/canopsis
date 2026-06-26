<template>
  <c-layout class="pt-2">
    <div :style="labelStyles">
      <h4 class="text-subtitle-1 font-weight-bold">
        {{ label }}
      </h4>
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
