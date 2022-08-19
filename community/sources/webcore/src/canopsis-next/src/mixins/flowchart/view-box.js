export const viewBoxMixin = {
  props: {
    viewBox: {
      type: Object,
      required: true,
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
  mounted() {
    this.$refs.svg.addEventListener('wheel', this.moveViewBoxByWheelEvent);

    document.addEventListener('transitionend', this.setViewBoxAfterTransition);
  },
  beforeDestroy() {
    this.$refs.svg.removeEventListener('wheel', this.moveViewBoxByWheelEvent);

    document.removeEventListener('transitionend', this.setViewBoxAfterTransition);
  },
  methods: {
    updateViewBox() {
      this.$emit('update:viewBox', this.viewBoxObject);
    },

    setViewBox(event) {
      const { width, height } = this.$refs.svg.getBoundingClientRect();

      if (event) {
        const widthDiff = (this.editorSize.width - width) * this.widthScale;
        const heightDiff = (this.editorSize.height - height) * this.heightScale;

        this.viewBoxObject.width -= widthDiff;
        this.viewBoxObject.height -= heightDiff;
      } else {
        this.viewBoxObject.width = width;
        this.viewBoxObject.height = height;
      }

      this.updateViewBox();

      this.editorSize.width = width;
      this.editorSize.height = height;
    },

    setViewBoxAfterTransition(event) {
      const isContainsSvg = event.target.contains(this.$refs.svg);

      if (isContainsSvg && event.propertyName === 'transform') {
        this.setViewBox();

        document.removeEventListener('transitionend', this.setViewBoxAfterTransition);
      }
    },

    moveViewBoxByWheelEvent(event) {
      event.preventDefault();

      const delta = event.deltaY;

      if (event.ctrlKey) {
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
        this.viewBoxObject.x += delta;
      }

      if (event.altKey) {
        this.viewBoxObject.y += delta;
      }
    },

    moveViewBox({ x, y }) {
      this.viewBoxObject.x -= x * this.widthScale;
      this.viewBoxObject.y -= y * this.heightScale;
    },
  },
};
