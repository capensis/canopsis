<template>
  <g>
    <path
      v-for="connectorTriangle in connectorTriangles"
      :key="connectorTriangle.side"
      :d="connectorTriangle.path"
      fill="transparent"
      @mouseenter="onMouseEnter(connectorTriangle.side)"
      @mouseleave="onMouseLeave"
      @mouseup="onConnectFinish(connectorTriangle)"
      @mousemove.stop=""
    />
  </g>
</template>

<script>
import { CONNECTOR_SIDES } from '@/constants';

import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';

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
      default: 5,
    },
  },
  data() {
    return {
      activeSide: undefined,
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
      return this.leftX + this.width / 2;
    },

    centerY() {
      return this.topY + this.height / 2;
    },

    topTrianglePath() {
      return [
        `M ${this.leftX} ${this.topY}`,
        `L ${this.rightX} ${this.topY}`,
        `L ${this.centerX} ${this.centerY}`,
      ].join('');
    },

    rightTrianglePath() {
      return [
        `M ${this.rightX} ${this.topY}`,
        `L ${this.rightX} ${this.bottomY}`,
        `L ${this.centerX} ${this.centerY}`,
      ].join('');
    },

    bottomTrianglePath() {
      return [
        `M ${this.rightX} ${this.bottomY}`,
        `L ${this.leftX} ${this.bottomY}`,
        `L ${this.centerX} ${this.centerY}`,
      ].join('');
    },

    leftTrianglePath() {
      return [
        `M ${this.leftX} ${this.bottomY}`,
        `L ${this.leftX} ${this.topY}`,
        `L ${this.centerX} ${this.centerY}`,
      ].join('');
    },

    connectorTriangles() {
      return [
        {
          path: this.topTrianglePath,
          side: CONNECTOR_SIDES.top,
        },
        {
          path: this.rightTrianglePath,
          side: CONNECTOR_SIDES.right,
        },
        {
          path: this.bottomTrianglePath,
          side: CONNECTOR_SIDES.bottom,
        },
        {
          path: this.leftTrianglePath,
          side: CONNECTOR_SIDES.left,
        },
      ];
    },
  },
  methods: {
    onMouseEnter(side) {
      this.activeSide = side;

      const point = calculateConnectorPointBySide(
        { x: this.x, y: this.y, width: this.width, height: this.height },
        side,
      );

      this.$emit('connecting', point);
    },

    onMouseLeave() {
      this.activeSide = undefined;

      this.$emit('unconnect');
    },

    onConnectFinish() {
      this.$emit('connected', { side: this.activeSide });
    },
  },
};
</script>
