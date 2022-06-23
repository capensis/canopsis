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
import { CONNECTOR_SIDES } from '@/plugins/flowchart/constants';

import { calculateConnectorPointBySide } from '@/plugins/flowchart/utils/connectors';

import PointsPath from '@/plugins/flowchart/components/common/points-path.vue';

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
  },
  data() {
    return {
      activeSide: undefined,
    };
  },
  computed: {
    rightX() {
      return this.x + this.width;
    },

    bottomY() {
      return this.y + this.height;
    },

    selectionWidth() {
      return this.rightX - this.x;
    },

    selectionHeight() {
      return this.bottomY - this.y;
    },

    centerX() {
      return this.x + this.width / 2;
    },

    centerY() {
      return this.y + this.height / 2;
    },

    centerPoint() {
      return { x: this.centerX, y: this.centerY };
    },

    topLeftPoint() {
      return { x: this.x, y: this.y };
    },

    topRightPoint() {
      return { x: this.rightX, y: this.y };
    },

    bottomRightPoint() {
      return { x: this.rightX, y: this.bottomY };
    },

    bottomLeftPoint() {
      return { x: this.x, y: this.bottomY };
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
