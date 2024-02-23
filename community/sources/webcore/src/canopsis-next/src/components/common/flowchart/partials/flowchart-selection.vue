<template>
  <path
    :d="path"
    :fill="color"
    :stroke="color"
    fill-opacity="0.1"
    stroke-width="1"
  />
</template>

<script>
export default {
  props: {
    color: {
      type: String,
      default: 'blue',
    },
    start: {
      type: Object,
      required: true,
    },
    end: {
      type: Object,
      required: true,
    },
  },
  computed: {
    path() {
      return [
        `M ${this.start.x} ${this.start.y}`,
        `L ${this.start.x} ${this.end.y}`,
        `L ${this.end.x} ${this.end.y}`,
        `L ${this.end.x} ${this.start.y}`,
        'Z',
      ].join('');
    },
  },
};
</script>
