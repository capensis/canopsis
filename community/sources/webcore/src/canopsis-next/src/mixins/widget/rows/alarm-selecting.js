import { throttle } from 'lodash';

export const widgetRowsSelectingAlarmMixin = {
  props: {
    selectable: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      prevEvent: null,
      selecting: false,
      selected: [],
    };
  },

  computed: {
    selectingMousePressed() {
      return this.selecting && !!this.prevEvent;
    },
  },

  created() {
    this.throttledMousemoveHandler = throttle(this.mousemoveHandler, 50);
  },

  async mounted() {
    if (this.selectable) {
      window.addEventListener('keydown', this.enableSelecting);
      window.addEventListener('keyup', this.disableSelecting);
      window.addEventListener('mousedown', this.mousedownHandler);
      window.addEventListener('mouseup', this.mouseupHandler);
    }
  },
  updated() {
    if (this.selecting) {
      this.calculateRowsPositions();
    }
  },
  beforeDestroy() {
    window.removeEventListener('keydown', this.enableSelecting);
    window.removeEventListener('keyup', this.disableSelecting);
    window.removeEventListener('mousedown', this.mousedownHandler);
    window.removeEventListener('mouseup', this.mouseupHandler);
  },

  methods: {
    calculateRowsPositions() {
      this.rowsPositions = Object.entries(this.$refs).reduce((acc, [key, value]) => {
        if (!key.startsWith('row') || !value) {
          return acc;
        }

        const position = value.$el.getBoundingClientRect();

        acc.push({
          position: {
            x1: position.x,
            x2: position.x + position.width,
            y1: position.y,
            y2: position.y + position.height,
          },
          alarm: value.$options.propsData.alarm,
        });

        return acc;
      }, []);
    },

    getIntersectAlarmsByPosition(newX, newY, prevX, prevY) {
      return this.rowsPositions?.reduce((acc, { position, alarm }) => {
        if (
          (prevX >= position.x1 && prevX <= position.x2 && prevY >= position.y1 && prevY <= position.y2)
          || (newX < position.x1 && prevX < position.x1)
          || (newX > position.x2 && prevX > position.x2)
          || (newY < position.y1 && prevY < position.y1)
          || (newY > position.y2 && prevY > position.y2)
        ) {
          return acc;
        }

        acc.push(alarm);

        return acc;
      }, []) ?? [];
    },

    mousedownHandler(event) {
      this.prevEvent = event;
    },

    mouseupHandler() {
      this.prevEvent = null;
    },

    mousemoveHandler(event) {
      if (!event.ctrlKey || !event.buttons || !this.prevEvent) {
        return;
      }

      const alarms = this.getIntersectAlarmsByPosition(
        event.clientX,
        event.clientY,
        this.prevEvent.clientX,
        this.prevEvent.clientY,
      );

      this.prevEvent = event;

      alarms.forEach(alarm => this.toggleSelected(alarm));
    },

    toggleSelected(alarm) {
      const index = this.selected.findIndex(({ _id: id }) => id === alarm._id);

      if (index === -1) {
        this.selected.push(alarm);

        return;
      }

      this.selected.splice(index, 1);
    },

    clearSelected() {
      this.selected = [];
    },

    enableSelecting({ key }) {
      if (key === 'Control') {
        this.selecting = true;
      }
    },

    disableSelecting({ key }) {
      if (key === 'Control') {
        this.selecting = false;
      }
    },
  },
};
