import { get } from 'lodash';

export default {
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && get(this.calendarEvent, 'data.meta.hasPopover');
    },

    ending() {
      return this.isPlaceholderWithDay
        ? this.isPlaceholderWithDay.sameDay(this.calendarEvent.end)
        : this.calendarEvent.ending || (
          this.calendar &&
        this.calendar.filled.end.sameDay(this.calendarEvent.end)
        );
    },

    canResize() {
      return !this.$dayspan.readOnly && this.ending;
    },
  },
  beforeDestroy() {
    this.resizeEndHandler();
  },
  methods: {
    resizeStartHandler(event) {
      if (event.button === 0) {
        event.stopPropagation();
        this.$emit('mouse-start-resize', event, this.calendarEvent);
        document.addEventListener('mouseup', this.resizeEndHandler);
      }
    },

    handlesEvents(event) {
      if (event) {
        event.stopPropagation();
      }

      return !!event;
    },

    resizeEndHandler() {
      document.removeEventListener('mouseup', this.resizeEndHandler);
    },

    mouseDownEvent($event) {
      if (this.handlesEvents($event)) {
        this.$emit('mouse-down-event', this.getEvent('mouse-down-event', $event));
      }
    },

    editCheck($event) {
      if (this.handlesEvents($event)) {
        this.menu = !this.menu;
      }
    },
  },
};
