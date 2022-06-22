<template lang="pug">
  g
    template(v-if="resizable")
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
      circle(
        v-for="circle in resizeCircles",
        :cx="circle.x",
        :cy="circle.y",
        :fill="color",
        :r="cornerRadius",
        :cursor="`${circle.direction}-resize`",
        @mousedown.stop="startResize(circle.direction)"
      )
    template(v-if="connectable")
      circle(
        v-for="connector in connectors",
        :key="connector._id",
        :cx="connector.x",
        :cy="connector.y",
        :fill="color",
        :r="cornerRadius"
      )
</template>

<script>
import { DIRECTIONS } from '@/plugins/flowchart/constants';

export default {
  props: {
    x: {
      type: Number,
      required: true,
    },
    y: {
      type: Number,
      required: true,
    },
    width: {
      type: Number,
      required: true,
    },
    height: {
      type: Number,
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
    connectable: {
      type: Boolean,
      default: false,
    },
    connectors: {
      type: Array,
      default: () => [],
    },
  },
  computed: {
    leftX() {
      return this.x - this.padding;
    },

    topY() {
      return this.y - this.padding;
    },

    rightX() {
      return this.x + this.width + this.padding;
    },

    bottomY() {
      return this.y + this.height + this.padding;
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

    resizeCircles() {
      return [
        { x: this.leftX, y: this.topY, direction: DIRECTIONS.northWest },
        { x: this.centerX, y: this.topY, direction: DIRECTIONS.north },
        { x: this.rightX, y: this.topY, direction: DIRECTIONS.northEast },
        { x: this.rightX, y: this.centerY, direction: DIRECTIONS.east },
        { x: this.rightX, y: this.bottomY, direction: DIRECTIONS.southEast },
        { x: this.centerX, y: this.bottomY, direction: DIRECTIONS.south },
        { x: this.leftX, y: this.bottomY, direction: DIRECTIONS.southWest },
        { x: this.leftX, y: this.centerY, direction: DIRECTIONS.west },
      ];
    },
  },
  methods: {
    startResize(direction) {
      this.$emit('start:resize', direction);
    },
  },
};
</script>
