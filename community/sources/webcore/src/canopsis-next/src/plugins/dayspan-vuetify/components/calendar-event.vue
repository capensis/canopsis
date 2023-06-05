<template lang="pug">
  v-menu(
    :content-class="contentClass",
    :disabled="!hasPopover",
    v-model="menu",
    v-bind="popoverProps"
  )
    template(#activator="{ on }")
      div.ds-calendar-event-menu(:class="classWithKey")
        div.ds-calendar-event(
          :style="style",
          @click.stop="editCheck",
          @mouseenter="mouseEnterEvent",
          @mouseleave="mouseLeaveEvent",
          @mousedown="mouseDownEvent",
          @mouseup="mouseUpEvent"
        )
          span(v-if="showName")
            slot(name="eventTitle", v-bind="{ calendarEvent, hasPrefix, getPrefix, details }")
              v-icon.ds-ev-icon.pr-1(
                v-if="hasIcon",
                size="14",
                :style="{ color: details.forecolor }"
              ) {{ details.icon }}
              span(v-if="hasPrefix") {{ getPrefix }}
              strong.ds-ev-title {{ details.title }}
              span.ds-ev-description {{ details.description }}
          span(v-else)
            slot(name="eventEmpty", v-bind="{ calendarEvent, details }") &nbsp;
          div.ds-calendar-event-resize(v-show="canResize", @mousedown="resizeStartHandler")
    slot(
      name="eventPopover",
      v-if="isShownPopover",
      v-bind="{ calendarEvent, calendar, edit, details, close: closePopover }"
    )
</template>

<script>
import { Functions as fn } from 'dayspan';
import { DsCalendarEvent } from 'dayspan-vuetify/src/components';

import eventMixin from '../mixins/event';

export default {
  extends: DsCalendarEvent,
  mixins: [eventMixin],
  computed: {
    fullStyles() {
      return this.$dayspan.getStyleFull(this.details, this.calendarEvent, this.index);
    },

    placeholderFullStyles() {
      const stateColor = this.$dayspan.getStyleColor(this.details, this.calendarEvent);
      const styles = this.$dayspan.getStylePlaceholderFull(
        this.details,
        this.calendarEvent,
        this.index,
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
        $event,
        type,
        calendarEvent: this.calendarEvent,
        closePopover: this.closePopover,
        openPopover: this.openPopover,
        calendar: this.calendar,
        details: this.details,
        day: this.isPlaceholderWithDay || this.$parent.day,
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
  .ds-calendar-event-menu {
    position: relative;
  }

  .ds-calendar-event {
    box-sizing: border-box;
    margin: 0 !important;
  }

  .ds-calendar-event-resize {
    pointer-events: auto;
    content: '';
    height: 100%;
    cursor: ew-resize;
    position: absolute;
    right: 0;
    width: 10px;
    top: 0;
    bottom: 0;
  }
</style>
