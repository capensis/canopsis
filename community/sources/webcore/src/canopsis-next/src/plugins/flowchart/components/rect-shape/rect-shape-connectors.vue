<template lang="pug">
  g
    circle(
      v-for="connector in connectors",
      :cx="connector.x",
      :cy="connector.y",
      :fill="color",
      :r="cornerRadius"
    )
</template>

<script>
import { range } from 'lodash';

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
    rx: {
      type: Number,
      default: 0,
    },
    ry: {
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
    connectorCount: {
      type: Number,
      default: 3,
    },
  },
  computed: {
    startX() {
      return this.x + this.rx;
    },

    widthConnectorOffset() {
      return (this.width - this.rx * 2) / (this.connectorCount + 1);
    },

    heightConnectorOffset() {
      return this.height / (this.connectorCount + 1);
    },

    sideConnectors() {
      return range(1, this.connectorCount + 1);
    },

    topConnectors() {
      return this.sideConnectors.map(index => ({
        x: this.startX + this.widthConnectorOffset * index,
        y: this.y,
      }));
    },

    bottomConnectors() {
      return this.sideConnectors.map(index => ({
        x: this.startX + this.widthConnectorOffset * index,
        y: this.y + this.height,
      }));
    },

    rightConnectors() {
      return this.sideConnectors.map(index => ({
        x: this.x + this.width,
        y: this.y + this.heightConnectorOffset * index,
      }));
    },

    leftConnectors() {
      return this.sideConnectors.map(index => ({
        x: this.x,
        y: this.y + this.heightConnectorOffset * index,
      }));
    },

    anglesConnectors() {
      return [
        { x: this.x, y: this.y },
        { x: this.x + this.width, y: this.y },
        { x: this.x + this.width, y: this.y + this.height },
        { x: this.x, y: this.y + this.height },
      ];
    },

    connectors() {
      const connectors = [
        ...this.topConnectors,
        ...this.rightConnectors,
        ...this.bottomConnectors,
        ...this.leftConnectors,
      ];

      if (!this.rx) {
        connectors.push(...this.anglesConnectors);
      }

      return connectors;
    },
  },
};
</script>
