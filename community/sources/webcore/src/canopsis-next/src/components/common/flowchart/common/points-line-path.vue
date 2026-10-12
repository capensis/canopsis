<template>
  <path
    ref="path"
    :d="path"
    :fill="fill"
    pointer-events="stroke"
    v-on="$listeners"
  />
</template>

<script>
import { LINE_TYPES } from '@/constants';

import { calculateCenterBetweenPoint } from '@/helpers/flowchart/points';

export default {
  props: {
    points: {
      type: Array,
      required: true,
    },
    fill: {
      type: String,
      default: 'transparent',
    },
    type: {
      type: String,
      default: LINE_TYPES.line,
    },
    curveOffset: {
      type: Number,
      default: 0.60,
    },
  },
  computed: {
    start() {
      return this.points[0];
    },

    end() {
      return this.points[this.points.length - 1];
    },

    diffX() {
      return this.end.x - this.start.x;
    },

    diffY() {
      return this.end.y - this.start.y;
    },

    curvePixelOffsetX() {
      return this.diffX * this.curveOffset;
    },

    curvePixelOffsetY() {
      return this.diffY * this.curveOffset;
    },

    curvesControls() {
      if (this.type === LINE_TYPES.verticalCurve) {
        return [
          { x: this.start.x, y: this.end.y - this.curvePixelOffsetY },
          { x: this.end.x, y: this.start.y + this.curvePixelOffsetY },
        ];
      }

      return [
        { x: this.end.x - this.curvePixelOffsetX, y: this.start.y },
        { x: this.start.x + this.curvePixelOffsetX, y: this.end.y },
      ];
    },

    startMoveTo() {
      return `M ${this.start.x} ${this.start.y}`;
    },

    linePoints() {
      return [
        this.startMoveTo,
        `L ${this.end.x} ${this.end.y}`,
      ];
    },

    curvesPoints() {
      const [firstControl, secondControl] = this.curvesControls;

      const centerPoint = calculateCenterBetweenPoint(firstControl, secondControl);

      return [
        this.startMoveTo,
        `Q ${firstControl.x} ${firstControl.y} ${centerPoint.x} ${centerPoint.y}`,
        `Q ${secondControl.x} ${secondControl.y} ${this.end.x} ${this.end.y}`,
      ];
    },

    leftCurvedAnglePoints() {
      const corner = { x: this.start.x, y: this.end.y };

      return [
        this.startMoveTo,
        `Q ${corner.x} ${corner.y} ${this.end.x} ${this.end.y}`,
      ];
    },

    rightCurvedAnglePoints() {
      const corner = { x: this.end.x, y: this.start.y };

      return [
        this.startMoveTo,
        `Q ${corner.x} ${corner.y} ${this.end.x} ${this.end.y}`,
      ];
    },

    rightElbowPoints() {
      return [
        this.startMoveTo,
        `h ${this.diffX}`,
        `v ${this.diffY}`,
      ];
    },

    leftElbowPoints() {
      return [
        this.startMoveTo,
        `v ${this.diffY}`,
        `h ${this.diffX}`,
      ];
    },

    horizontalSteppedPoints() {
      const halfDiffX = this.diffX / 2;

      return [
        this.startMoveTo,
        `h ${halfDiffX}`,
        `v ${this.diffY}`,
        `h ${halfDiffX}`,
      ];
    },

    verticalSteppedPoints() {
      const halfDiffY = this.diffY / 2;

      return [
        this.startMoveTo,
        `v ${halfDiffY}`,
        `h ${this.diffX}`,
        `v ${halfDiffY}`,
      ];
    },

    path() {
      const points = {
        [LINE_TYPES.line]: this.linePoints,
        [LINE_TYPES.horizontalCurve]: this.curvesPoints,
        [LINE_TYPES.verticalCurve]: this.curvesPoints,
        [LINE_TYPES.leftCurvedAngle]: this.leftCurvedAnglePoints,
        [LINE_TYPES.rightCurvedAngle]: this.rightCurvedAnglePoints,
        [LINE_TYPES.rightElbow]: this.rightElbowPoints,
        [LINE_TYPES.leftElbow]: this.leftElbowPoints,
        [LINE_TYPES.horizontalStepped]: this.horizontalSteppedPoints,
        [LINE_TYPES.verticalStepped]: this.verticalSteppedPoints,
      }[this.type];

      return points.join('');
    },
  },
  methods: {
    getCenterPoint() {
      const length = this.$refs.path.getTotalLength();

      return this.$refs.path.getPointAtLength(length * 0.5);
    },
  },
};
</script>
