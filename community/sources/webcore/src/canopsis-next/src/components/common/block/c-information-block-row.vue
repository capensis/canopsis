<template lang="pug">
  v-layout.pt-2(row)
    v-flex(:style="labelStyles")
      h4.subheading.font-weight-bold {{ label }}
    v-flex
      slot {{ value }}
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
