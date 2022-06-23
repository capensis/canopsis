<template lang="pug">
  g
    points-path(
      v-for="connectorTriangle in connectorTriangles",
      :key="connectorTriangle.side",
      :points="connectorTriangle.points",
      @mouseenter="onMouseEnter(connectorTriangle.side)",
      @mouseleave="onMouseLeave",
      @mouseup="onConnectFinish(connectorTriangle)",
      @mousemove.stop=""
    )
</template>

<script>
import { CONNECTOR_SIDES } from '@/constants';

import { calculateConnectorPointBySide } from '@/helpers/flowchart/connectors';

import PointsPath from './points-path.vue';

export default {
  components: { PointsPath },
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

    centerPoint() {
      return { x: this.centerX, y: this.centerY };
    },

    topLeftPoint() {
      return { x: this.leftX, y: this.topY };
    },

    topRightPoint() {
      return { x: this.rightX, y: this.topY };
    },

    bottomRightPoint() {
      return { x: this.rightX, y: this.bottomY };
    },

    bottomLeftPoint() {
      return { x: this.leftX, y: this.bottomY };
    },

    topTrianglePoints() {
      return [
        this.topLeftPoint,
        this.topRightPoint,
        this.centerPoint,
      ];
    },

    rightTrianglePoints() {
      return [
        this.topRightPoint,
        this.bottomRightPoint,
        this.centerPoint,
      ];
    },

    bottomTrianglePoints() {
      return [
        this.bottomRightPoint,
        this.bottomLeftPoint,
        this.centerPoint,
      ];
    },

    leftTrianglePoints() {
      return [
        this.bottomLeftPoint,
        this.topLeftPoint,
        this.centerPoint,
      ];
    },

    connectorTriangles() {
      return [
        {
          points: this.topTrianglePoints,
          side: CONNECTOR_SIDES.top,
        },
        {
          points: this.rightTrianglePoints,
          side: CONNECTOR_SIDES.right,
        },
        {
          points: this.bottomTrianglePoints,
          side: CONNECTOR_SIDES.bottom,
        },
        {
          points: this.leftTrianglePoints,
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
