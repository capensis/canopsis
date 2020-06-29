<template lang="pug">
  v-menu(
    class="ds-calendar-event-menu",
    :content-class="contentClass",
    :disabled="!hasPopover",
    v-model="menu",
    v-bind="popoverProps"
  )
    .ds-calendar-event(
      slot="activator",
      :style="style",
      @click.stop="editCheck",
      @mouseenter="mouseEnterEvent",
      @mouseleave="mouseLeaveEvent",
      @mousedown="mouseDownEvent",
      @mouseup="mouseUpEvent"
    )
      span(v-if="showName")
        slot(name="eventTitle", v-bind="{ calendarEvent, hasPrefix, getPrefix, details }")
          v-icon.ds-ev-icon(
            v-if="hasIcon",
            size="14",
            :style="{ color: details.forecolor }"
          ) {{ details.icon }}
          span(v-if="hasPrefix") {{ getPrefix }}
          strong.ds-ev-title {{ details.title }}
          span.ds-ev-description {{ details.description }}

      span(v-else)
        slot(name="eventEmpty", v-bind="{ calendarEvent, details }") &nbsp;
      .ds-calendar-event-resize(v-show="canResize", @mousedown="resizeStartHandler")
    slot(name="eventPopover", v-if="menu", v-bind="{ calendarEvent, calendar, edit, details, close }")
</template>

<script>
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
    cursor: e-resize;
    position: absolute;
    right: 0;
    width: 10px;
    top: 0;
    bottom: 0;
  }
</style>
