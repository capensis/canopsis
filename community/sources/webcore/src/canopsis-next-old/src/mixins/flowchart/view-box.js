import { isEmpty } from 'lodash';

import { calculateCenterBetweenPoint } from '@/helpers/flowchart/points';
import { getShapesBounds } from '@/helpers/flowchart/shapes';

export const viewBoxMixin = {
  props: {
    viewBox: {
      type: Object,
      required: true,
    },
    fitPaddingPercent: {
      type: Number,
      default: 0.1,
    },
  },
  data() {
    return {
      viewBoxObject: this.viewBox,
      editorSize: {
        width: 0,
        height: 0,
      },
    };
  },
  computed: {
    viewBoxString() {
      const { x, y, width, height } = this.viewBoxObject;

      return `${x} ${y} ${width} ${height}`;
    },

    widthScale() {
      return this.viewBoxObject.width / this.editorSize.width;
    },

    heightScale() {
      return this.viewBoxObject.height / this.editorSize.height;
    },
  },
  created() {
    this.resizeObserver = new ResizeObserver(this.updateViewBoxByResize);
  },
  mounted() {
    this.setViewBox();
    this.fitToShapes();

    this.$flowchart.on('wheel', this.moveViewBoxByWheelEvent);

    this.resizeObserver.observe(this.$parent.$el);
  },
  beforeDestroy() {
    this.$flowchart.off('wheel', this.moveViewBoxByWheelEvent);

    this.resizeObserver.unobserve(this.$parent.$el);
    this.resizeObserver.disconnect();
  },
  methods: {
    updateViewBox() {
      this.$emit('update:viewBox', this.viewBoxObject);
    },

    updateEditorSize(width, height) {
      this.editorSize.width = width;
      this.editorSize.height = height;
    },

    getSvgSizes() {
      const style = window.getComputedStyle(this.$refs.svg);

      return {
        width: parseInt(style.width, 10),
        height: parseInt(style.height, 10),
      };
    },

    fitToShapes() {
      if (isEmpty(this.shapes)) {
        return;
      }

      const { min, max } = getShapesBounds(this.shapes);

      const paddingFactor = 1 + this.fitPaddingPercent;

      const optimalHeight = (max.y - min.y) * paddingFactor;
      const optimalWidth = (max.x - min.x) * paddingFactor;

      if (optimalWidth > this.viewBoxObject.width) {
        this.viewBoxObject.height *= optimalWidth / this.viewBoxObject.width;
        this.viewBoxObject.width = optimalWidth;
      }

      if (optimalHeight > this.viewBoxObject.height) {
        this.viewBoxObject.width *= optimalHeight / this.viewBoxObject.height;
        this.viewBoxObject.height = optimalHeight;
      }

      const center = calculateCenterBetweenPoint(max, min);

      this.viewBoxObject.x = center.x - this.viewBoxObject.width / 2;
      this.viewBoxObject.y = center.y - this.viewBoxObject.height / 2;
    },

    setViewBox() {
      const { width, height } = this.getSvgSizes();

      this.viewBoxObject.width = width;
      this.viewBoxObject.height = height;

      this.updateViewBox();
      this.updateEditorSize(width, height);
    },

    updateViewBoxByResize() {
      const { width, height } = this.getSvgSizes();

      const widthDiff = (this.editorSize.width - width) * this.widthScale;
      const heightDiff = (this.editorSize.height - height) * this.heightScale;

      this.viewBoxObject.width -= widthDiff;
      this.viewBoxObject.x += widthDiff / 2;
      this.viewBoxObject.height -= heightDiff;
      this.viewBoxObject.y += heightDiff / 2;

      this.updateEditorSize(width, height);
      this.updateViewBox();
    },

    moveViewBoxByWheelEvent({ event }) {
      const delta = event.deltaY;

      if (event.ctrlKey) {
        event.preventDefault();

        const percent = delta < 0 ? 0.05 : -0.05;

        const scaleWidth = this.viewBoxObject.width * percent;
        const scaleHeight = this.viewBoxObject.height * percent;

        const deltaWidth = scaleWidth * 2;
        const deltaHeight = scaleHeight * 2;

        const { x, y } = this.normalizeCursor({ x: event.clientX, y: event.clientY });

        const cursorPercentX = (x - this.viewBoxObject.x) / this.viewBoxObject.width;
        const cursorPercentY = (y - this.viewBoxObject.y) / this.viewBoxObject.height;

        const offsetX = deltaWidth * cursorPercentX;
        const offsetY = deltaHeight * cursorPercentY;
        const offsetWidth = scaleWidth + deltaWidth - offsetX;
        const offsetHeight = scaleHeight + deltaHeight - offsetY;

        this.viewBoxObject.x += offsetX;
        this.viewBoxObject.y += offsetY;
        this.viewBoxObject.width -= offsetWidth;
        this.viewBoxObject.height -= offsetHeight;
      }

      if (event.shiftKey) {
        event.preventDefault();

        this.viewBoxObject.x += delta;
      }

      if (event.altKey) {
        event.preventDefault();

        this.viewBoxObject.y += delta;
      }
    },

    moveViewBox({ x, y }) {
      this.viewBoxObject.x -= x * this.widthScale;
      this.viewBoxObject.y -= y * this.heightScale;
    },
  },
};
