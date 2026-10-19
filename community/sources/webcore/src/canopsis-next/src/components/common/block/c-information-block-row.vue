<template>
  <div class="d-flex align-center ga-4 pt-2">
    <div :style="labelStyles">
      <div class="font-weight-bold text-label-large">
        {{ label }}
      </div>
    </div>
    <div class="text-break">
      <slot>{{ value }}</slot>
    </div>
  </div>
</template>

<script>
import { computed } from 'vue';

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
  setup(props) {
    const labelStyles = computed(() => ({ minWidth: `${props.width}px`, maxWidth: `${props.width}px` }));

    return { labelStyles };
  },
};
</script>
