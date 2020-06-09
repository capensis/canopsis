<template lang="pug">
  v-menu.ds-calendar-event(
    :content-class="contentClass",
    :disabled="!hasPopover",
    :style="style",
    v-model="menu",
    v-bind="popoverProps"
  )
    template(slot="activator")
      .ds-calendar-event-span(
        @click="editCheck",
        @mouseenter="mouseEnterEvent",
        @mouseleave="mouseLeaveEvent",
        @mousedown="mouseDownEvent",
        @mouseup="mouseUpEvent"
      )
        span(v-if="showName")
          slot(name="eventTimeTitle", v-bind="{ calendarEvent, details }")
            v-icon.ds-ev-icon(
              v-if="hasIcon",
              size="14",
              :style="{ color: details.forecolor }"
            ) {{ details.icon }}
            strong.ds-ev-title {{ details.title }}
            span.ds-ev-description {{ details.description }}
        span(v-else)
          slot(name="eventTimeEmpty", v-bind="{ calendarEvent, details }")
      .ds-calendar-event-resize(v-show="calendarEvent.ending", @mousedown="resizeStartHandler")
    slot(name="eventPopover", v-bind="{ calendarEvent, calendar, edit, details, close }")
</template>

<script>
import { get } from 'lodash';
import { DsCalendarEventTime } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEventTime,
  computed: {
    hasPopover() {
      return !!this.$scopedSlots.eventPopover && get(this.calendarEvent, 'data.meta.hasPopover');
    },

    style() {
      return this.isPlaceholderWithDay ?
        this.$dayspan.getStylePlaceholderTimed(this.details, this.calendarEvent, this.isPlaceholderWithDay) :
        this.$dayspan.getStyleTimed(this.details, this.calendarEvent);
    },
  },
  methods: {
    resizeStartHandler(event) {
      event.stopPropagation();
      this.$emit('mouse-start-resize', event, this.calendarEvent);
      document.addEventListener('mouseup', this.resizeEndHandler);
    },
    resizeEndHandler() {
      document.removeEventListener('mouseup', this.resizeEndHandler);
    },
  },
};
</script>

<style lang="scss">
  .ds-calendar-event-resize {
    pointer-events: auto;
    content: '';
    width: 100%;
    cursor: n-resize;
    position: absolute;
    left: 0;
    right: 0;
    height: 10px;
    bottom: 0;
  }
</style>
