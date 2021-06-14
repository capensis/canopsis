import { getMenuClassByCalendarEvent } from '@/helpers/calendar/dayspan';

import popoverMixin from './popover';

export default {
  mixins: [popoverMixin],
  computed: {
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
      return !this.$dayspan.readOnly && this.ending && !this.calendarEvent.data.withoutResize;
    },

    classWithKey() {
      return getMenuClassByCalendarEvent(this.calendarEvent);
    },
  },
  beforeDestroy() {
    this.resizeEndHandler();
  },
  methods: {
    resizeStartHandler(event) {
      if (event.button === 0) {
        event.stopPropagation();
        this.$emit('mouse-start-resize', this.getEvent('mouse-start-resize', event, {
          calendarEvent: this.calendarEvent,
        }));
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
      this.$emit('edit-event', this.getEvent('mouse-down-event', {}, { calendarEvent }));
    },

    editCheck(event) {
      if (this.isPlaceholderWithDay) {
        this.handlesEvents(event);
        return;
      }

      if (!this.hasPopover) {
        this.$emit('edit', this.getEvent('edit', event));
        return;
      }

      if (this.handlesEvents(event) && !this.menu) {
        this.menu = true;

        if (!this.isPlaceholderWithDay) {
          this.$emit('mouse-start-edit', this.getEvent('mouse-start-edit', event));
          this.$emit('clear-placeholder');
        }
      }
    },

    closePopover(event) {
      this.menu = false;
      this.$emit('mouse-end-edit', this.getEvent('mouse-end-edit', event));
    },
  },
};
