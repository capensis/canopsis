<template lang="pug">
  v-menu(
    :content-class="contentClass",
    :disabled="!hasPopover",
    v-model="menu",
    v-bind="popoverProps"
  )
    template(#activator="{ on }")
      div.ds-calendar-event(v-on="on", :style="style", :class="classWithKey")
        div.ds-calendar-event-span(
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
        div.ds-calendar-event-time-resize(v-show="canResize", @mousedown="resizeStartHandler")
    slot(
      name="eventPopover",
      v-if="isShownPopover",
      v-bind="{ calendarEvent, calendar, edit, details, close: closePopover }"
    )
</template>

<script>
import { Functions as fn } from 'dayspan';
import { DsCalendarEventTime } from 'dayspan-vuetify/src/components';

import eventMixin from '../mixins/event';

export default {
  extends: DsCalendarEventTime,
  mixins: [eventMixin],
  computed: {
    fullStyles() {
      return this.$dayspan.getStyleTimed(this.details, this.calendarEvent);
    },

    placeholderFullStyles() {
      const stateColor = this.$dayspan.getStyleColor(this.details, this.calendarEvent);
      const styles = this.$dayspan.getStylePlaceholderTimed(
        this.details,
        this.calendarEvent,
        this.isPlaceholderWithDay,
      );
      styles.backgroundColor = this.$dayspan.blend(stateColor, 1.1, this.$dayspan.placeholderBlendTarget);

      return styles;
    },

    style() {
      return this.isPlaceholderWithDay ? this.placeholderFullStyles : this.fullStyles;
    },
  },
  methods: {
    getEvent(type, $event, extra = {}) {
      return fn.extend({
        type,
        $event,
        calendarEvent: this.calendarEvent,
        closePopover: this.closePopover,
        openPopover: this.openPopover,
        calendar: this.calendar,
        details: this.details,
        day: this.isPlaceholderWithDay || this.$parent.day,
        offset: this.getRelativeTimeAt($event),
        left: $event.button === 0,
        right: $event.button === 1,
        handled: false,
        $vm: this,
        $element: this.$el,
      }, extra);
    },
  },
};
</script>

<style lang="scss">
  .ds-calendar-event {
    .ds-day & {
      border: 1px solid white;
    }

    & > .v-menu__activator {
      height: 100%;
    }

    .ds-calendar-event-time-resize {
      pointer-events: auto;
      content: '';
      width: 100%;
      cursor: ns-resize;
      position: absolute;
      left: 0;
      right: 0;
      bottom: 0;
      height: 10px;
      max-height: 30%;
    }
  }
</style>
