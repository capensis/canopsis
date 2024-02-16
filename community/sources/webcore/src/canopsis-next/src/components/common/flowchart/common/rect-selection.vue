<template>
  <g>
    <rect
      :x="leftX"
      :y="topY"
      :width="selectionWidth"
      :height="selectionHeight"
      :stroke="color"
      fill="transparent"
      stroke-width="1"
      stroke-dasharray="4 4"
      pointer-events="none"
    />
    <circle
      v-for="circle in resizeCircles"
      :key="circle.direction"
      :cx="circle.x"
      :cy="circle.y"
      :fill="color"
      :r="cornerRadius"
      :cursor="`${circle.direction}-resize`"
      @mousedown.stop="startResize(circle.direction)"
    />
  </g>
</template>

<script>
import { DIRECTIONS } from '@/constants';

import { resizeRectangleShape } from '@/helpers/flowchart/resize';

export default {
  inject: ['$flowchart'],
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
    aspectRatio: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      direction: undefined,
      ratio: undefined,
    };
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
      this.$flowchart.on('mousemove', this.onResize);
      this.$flowchart.on('mouseup', this.finishResize);

      this.direction = direction;
      this.ratio = this.width !== 0 ? this.height / this.width : 1;
    },

    finishResize() {
      this.$flowchart.off('mousemove', this.onResize);
      this.$flowchart.off('mouseup', this.finishResize);
    },

    normalizeCursorByDirection(cursor) {
      const newCursor = {
        x: cursor.x,
        y: cursor.y,
      };

      if (this.direction.includes(DIRECTIONS.north)) {
        newCursor.y += this.padding;
      } else if (this.direction.includes(DIRECTIONS.south)) {
        newCursor.y -= this.padding;
      }

      if (this.direction.includes(DIRECTIONS.west)) {
        newCursor.x += this.padding;
      } else if (this.direction.includes(DIRECTIONS.east)) {
        newCursor.x -= this.padding;
      }

      return newCursor;
    },

    onResize({ event, cursor }) {
      const normalizedCursor = this.normalizeCursorByDirection(cursor);

      const rect = resizeRectangleShape({
        rect: {
          x: this.x,
          y: this.y,
          width: this.width,
          height: this.height,
        },
        cursorX: normalizedCursor.x,
        cursorY: normalizedCursor.y,
        direction: this.direction,
        ratio: this.ratio,
        aspectRatio: this.aspectRatio || event.shiftKey,
      });

      this.$emit('update', rect);
    },
  },
};
</script>
