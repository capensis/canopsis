import { get } from 'lodash';

export default {
  computed: {
    hasPopover() {
      return get(this.calendarEvent, 'data.meta.hasPopover', !!this.$scopedSlots.eventPopover && !this.isPlaceholderWithDay);
    },

    isPlaceholderSameDay() {
      return this.isPlaceholderWithDay.sameDay(this.calendarEvent.end);
    },

    isCalendarSameDay() {
      return this.calendar && this.calendar.filled.end.sameDay(this.calendarEvent.end);
    },

    ending() {
      return this.isPlaceholderWithDay
        ? this.isPlaceholderSameDay
        : this.calendarEvent.ending || this.isCalendarSameDay;
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

    mouseDownEvent(event) {
      if (this.handlesEvents(event)) {
        this.$emit('mouse-down-event', this.getEvent('mouse-down-event', event));
      }
    },

    edit(calendarEvent) {
      this.$emit('edit-event', calendarEvent);
      this.menu = false;
    },

    editCheck(event) {
      if (this.handlesEvents(event)) {
        this.menu = !this.menu;

        if (!this.isPlaceholderWithDay) {
          this.$emit('clear-placeholder');
        }
      }
    },
  },
};
