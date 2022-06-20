<template lang="pug">
  g
    rect(
      :x="leftX",
      :y="topY",
      :width="selectionWidth",
      :height="selectionHeight",
      :stroke="color",
      fill="transparent",
      stroke-width="1",
      stroke-dasharray="4 4",
      pointer-events="none"
    )
    template(v-if="resizable")
      // Top left corner
      circle(
        :cx="leftX",
        :cy="topY",
        :fill="color",
        :r="cornerRadius",
        cursor="nw-resize",
        @mousedown.stop="startResize('nw')"
      )
      // Top center
      circle(
        :cx="centerX",
        :cy="topY",
        :fill="color",
        :r="cornerRadius",
        cursor="n-resize",
        @mousedown.stop="startResize('n')"
      )
      // Top right corner
      circle(
        :cx="rightX",
        :cy="topY",
        :fill="color",
        :r="cornerRadius",
        cursor="ne-resize",
        @mousedown.stop="startResize('ne')"
      )
      // Right center
      circle(
        :cx="rightX",
        :cy="centerY",
        :fill="color",
        :r="cornerRadius",
        cursor="e-resize",
        @mousedown.stop="startResize('e')"
      )
      // Bottom right corner
      circle(
        :cx="rightX",
        :cy="bottomY",
        :fill="color",
        :r="cornerRadius",
        cursor="se-resize",
        @mousedown.stop="startResize('se')"
      )
      // Bottom center
      circle(
        :cx="centerX",
        :cy="bottomY",
        :fill="color",
        :r="cornerRadius",
        cursor="s-resize",
        @mousedown.stop="startResize('s')"
      )
      // Bottom left corner
      circle(
        :cx="leftX",
        :cy="bottomY",
        :fill="color",
        :r="cornerRadius",
        cursor="sw-resize",
        @mousedown.stop="startResize('sw')"
      )
      // Left center
      circle(
        :cx="leftX",
        :cy="centerY",
        :fill="color",
        :r="cornerRadius",
        cursor="w-resize",
        @mousedown.stop="startResize('w')"
      )
</template>

<script>
export default {
  props: {
    rect: {
      type: Object,
      required: true,
    },
    padding: {
      type: Number,
      default: 0,
    },
    color: {
      type: String,
      default: 'blue',
    },
    cornerRadius: {
      type: Number,
      default: 4,
    },
    resizable: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    leftX() {
      return this.rect.x - this.padding;
    },

    topY() {
      return this.rect.y - this.padding;
    },

    rightX() {
      return this.rect.x + this.rect.width + this.padding;
    },

    bottomY() {
      return this.rect.y + this.rect.height + this.padding;
    },

    selectionWidth() {
      return this.rightX - this.leftX;
    },

    selectionHeight() {
      return this.bottomY - this.topY;
    },

    centerX() {
      return this.leftX + this.selectionWidth / 2;
    },

    centerY() {
      return this.topY + this.selectionHeight / 2;
    },
  },
  methods: {
    startResize(direction) {
      this.$emit('start:resize', direction);
    },
  },
};
</script>
