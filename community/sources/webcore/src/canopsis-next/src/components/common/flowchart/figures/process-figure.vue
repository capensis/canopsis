<template lang="pug">
  g
    rect(
      v-on="$listeners",
      :x="x",
      :y="y",
      :rx="offset",
      :ry="offset",
      :width="width",
      :height="height",
      :stroke="stroke",
      :stroke-width="strokeWidth"
    )
    rect(
      v-on="$listeners",
      :x="insetX",
      :y="y",
      :width="insetWidth",
      :height="height",
      :stroke="stroke",
      :stroke-width="strokeWidth"
    )
</template>

<script>
export default {
  props: {
    width: {
      type: Number,
      required: true,
    },
    height: {
      type: Number,
      required: true,
    },
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    offset: {
      type: Number,
      default: 0,
    },
    strokeWidth: {
      type: Number,
      default: 0,
    },
    stroke: {
      type: String,
      required: false,
    },
  },
  computed: {
    halfWidth() {
      return this.width / 2;
    },

    insetX() {
      const x = this.x + this.offset + this.strokeWidth;
      const centerX = this.x + this.halfWidth;

      return x <= centerX ? x : this.x + this.halfWidth;
    },

    insetWidth() {
      const width = this.width - (this.offset + this.strokeWidth) * 2;

      return width > 0 ? width : 0;
    },
  },
};
</script>
