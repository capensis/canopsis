import { cloneDeep } from 'lodash';

import { calculateShapeIconPosition } from '@/helpers/flowchart/shapes';

export const mapFlowchartPointsMixin = {
  props: {
    points: {
      type: Array,
      default: () => [],
    },
    shapes: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pointsData: [],
    };
  },
  watch: {
    points: {
      immediate: true,
      handler(value) {
        this.pointsData = cloneDeep(value);
      },
    },
  },
  computed: {
    shapesIcons() {
      return this.pointsData.reduce((acc, point) => {
        const { shape: shapeId } = point;

        if (shapeId) {
          const { x, y } = calculateShapeIconPosition(this.shapes[shapeId]);

          acc.push({
            x: x - this.iconSize / 2,
            y,
            point,
          });
        }

        return acc;
      }, []);
    },

    nonShapesIcons() {
      return this.pointsData.reduce((acc, point) => {
        if (!point.shape) {
          acc.push({
            x: this.calculatePointX(point),
            y: this.calculatePointY(point),
            point,
          });
        }

        return acc;
      }, []);
    },
  },
  methods: {
    calculatePointX(point) {
      return point.x - this.iconSize / 2;
    },

    calculatePointY(point) {
      if (point.entity) {
        return point.y - this.iconSize;
      }

      return point.y - this.iconSize / 2;
    },
  },
};
