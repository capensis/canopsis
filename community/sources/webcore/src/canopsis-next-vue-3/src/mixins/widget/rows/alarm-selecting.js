import { throttle } from 'lodash';

import { hasAlarmSelectingModifier, isAlarmSelectingModifierKey } from '@/helpers/event';

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
    };
  },

  computed: {
    selectingMousePressed() {
      return this.selecting && !!this.prevEvent;
    },

    wrapperListeners() {
      return this.selectable ? { mousemove: this.throttledMousemoveHandler } : {};
    },
  },

  created() {
    this.throttledMousemoveHandler = throttle(this.mousemoveHandler, 50);
  },

  mounted() {
    if (this.selectable) {
      window.addEventListener('keydown', this.handleKeyDown);
      window.addEventListener('keyup', this.handleKeyUp);
      window.addEventListener('mousedown', this.mousedownHandler);
      window.addEventListener('mouseup', this.mouseupHandler);
      window.addEventListener('resize', this.handleSelectingResize);
      window.addEventListener('blur', this.disableSelecting);
    }
  },

  updated() {
    if (this.selecting) {
      this.calculateRowsPositions();
    }
  },

  beforeUnmount() {
    window.removeEventListener('keydown', this.handleKeyDown);
    window.removeEventListener('keyup', this.handleKeyUp);
    window.removeEventListener('mousedown', this.mousedownHandler);
    window.removeEventListener('mouseup', this.mouseupHandler);
    window.removeEventListener('resize', this.handleSelectingResize);
    window.removeEventListener('blur', this.disableSelecting);
  },

  methods: {
    calculateRowsPositions() {
      this.rowsPositions = Object.entries(this.$refs).reduce((acc, [key, value]) => {
        if (!key.startsWith('row') || !value) {
          return acc;
        }

        const instance = Array.isArray(value) ? value[0] : value;

        if (!instance?.$el) {
          return acc;
        }

        const position = instance.$el.getBoundingClientRect();

        acc.push({
          position: {
            x1: position.x,
            x2: position.x + position.width,
            y1: position.y,
            y2: position.y + position.height,
          },
          alarm: instance.alarm,
        });

        return acc;
      }, []);
    },

    getIntersectAlarmsByPosition(newX, newY, prevX, prevY) {
      return (
        this.rowsPositions?.reduce((acc, { position, alarm }) => {
          if (
            (prevX >= position.x1 &&
              prevX <= position.x2 &&
              prevY >= position.y1 &&
              prevY <= position.y2) ||
            (newX < position.x1 && prevX < position.x1) ||
            (newX > position.x2 && prevX > position.x2) ||
            (newY < position.y1 && prevY < position.y1) ||
            (newY > position.y2 && prevY > position.y2)
          ) {
            return acc;
          }

          acc.push(alarm);

          return acc;
        }, []) ?? []
      );
    },

    mousedownHandler(event) {
      this.prevEvent = event;
    },

    mouseupHandler() {
      this.prevEvent = null;
    },

    mousemoveHandler(event) {
      if (!hasAlarmSelectingModifier(event) || event.buttons !== 1 || !this.prevEvent) {
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
      this.toggleAlarmSelection(alarm, !this.isAlarmSelected(alarm));
    },

    enableSelecting() {
      this.selecting = true;
      this.calculateRowsPositions();
    },

    disableSelecting() {
      this.selecting = false;
    },

    handleSelectingResize() {
      if (this.selecting) {
        this.calculateRowsPositions();
      }
    },

    handleKeyDown({ key }) {
      if (isAlarmSelectingModifierKey(key)) {
        this.enableSelecting();
      }
    },

    handleKeyUp({ key }) {
      if (isAlarmSelectingModifierKey(key)) {
        this.disableSelecting();
      }
    },
  },
};
