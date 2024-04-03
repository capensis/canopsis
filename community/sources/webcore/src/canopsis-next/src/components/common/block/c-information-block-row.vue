<template>
  <v-layout class="pt-2">
    <v-flex :style="labelStyles">
      <h4 class="text-subtitle-1 font-weight-bold">
        {{ label }}
      </h4>
    </v-flex>
    <v-flex>
      <slot>{{ value }}</slot>
    </v-flex>
  </v-layout>
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
