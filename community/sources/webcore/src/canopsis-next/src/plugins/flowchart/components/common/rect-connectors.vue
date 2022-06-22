<template lang="pug">
  g
    rect(
      v-for="connectorRect in connectorRects",
      :key="connectorRect.side",
      :x="connectorRect.x",
      :y="connectorRect.y",
      :width="connectorRect.width",
      :height="connectorRect.height",
      :fill="color",
      :opacity="activeSide === connectorRect.side ? 1 : 0",
      @mouseenter="onMouseEnter(connectorRect.side)",
      @mouseleave="onMouseLeave",
      @mouseup="onConnectFinish(connectorRect)",
      @mousemove.stop.passive="onConnectorMove"
    )
</template>

<script>
import { CONNECTOR_SIDES } from '@/plugins/flowchart/constants';
import { calculateConnectorPointBySide } from '@/plugins/flowchart/utils/connectors';

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
    color: {
      type: String,
      default: 'blue',
    },
    size: {
      type: Number,
      default: 5,
    },
  },
  data() {
    return {
      activeSide: undefined,
      offset: {
        x: 0.5,
        y: 0.5,
      },
    };
  },
  computed: {
    connectorRects() {
      return [
        {
          x: this.x,
          y: this.y - this.size,
          width: this.width,
          height: this.size,
          side: CONNECTOR_SIDES.top,
        },
        {
          x: this.x + this.width,
          y: this.y,
          width: this.size,
          height: this.height,
          side: CONNECTOR_SIDES.right,
        },
        {
          x: this.x,
          y: this.y + this.height,
          width: this.width,
          height: this.size,
          side: CONNECTOR_SIDES.bottom,
        },
        {
          x: this.x - this.size,
          y: this.y,
          width: this.size,
          height: this.height,
          side: CONNECTOR_SIDES.left,
        },
      ];
    },
  },
  methods: {
    onMouseEnter(side) {
      this.activeSide = side;
    },

    onMouseLeave() {
      this.activeSide = undefined;

      this.$emit('unconnect');
    },

    onConnectFinish() {
      this.$emit('connected', {
        side: this.activeSide,
        offset: this.offset,
      });
    },

    onConnectorMove(event) {
      const { x, y } = event.target.getBoundingClientRect();
      const relativeX = event.clientX - x;
      const relativeY = event.clientY - y;

      this.offset = {
        x: relativeX && relativeX / this.width,
        y: relativeY && relativeY / this.height,
      };

      const point = calculateConnectorPointBySide(
        { x: this.x, y: this.y, width: this.width, height: this.height },
        this.activeSide,
        this.offset,
      );

      this.$emit('connecting', point);
    },
  },
};
</script>
