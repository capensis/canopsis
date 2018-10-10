<template>
  <v-menu
    class="ds-calendar-event-menu"
    :content-class="contentClass"
    :disabled="!hasPopover"
    v-model="menu"
    v-bind="popoverProps">

    <div class="ds-calendar-event"
         slot="activator"
         :style="style"
         @click.stop="editCheck">

      <span v-if="showName" class="ds-calendar-event-name">
        <slot name="eventTitle" v-bind="{calendarEvent, hasPrefix, getPrefix, details}">
          <strong class="ds-ev-title">{{ details.title }}</strong>
        </slot>
      </span>

      <span v-else>
        <slot name="eventEmpty" v-bind="{calendarEvent, details}">&nbsp;</slot>
      </span>

    </div>

    <slot name="eventPopover" v-bind="{calendarEvent, calendar, edit, details, close}"></slot>

  </v-menu>
</template>

<script>
import { DsCalendarEvent } from 'dayspan-vuetify/src/components';

export default {
  extends: DsCalendarEvent,
  computed: {
    hasPrefix() {
      return false;
    },
  },
  methods: {
    editCheck() {

    },
  },
};
</script>

<style lang="scss">
  .ds-calendar-event-menu {
    position: absolute;
    left: 0;
    top: 0;
    width: 100%;
    height: 100%!important;
    padding: 4px;

    .ds-calendar-event {
      padding-left: 0;
      display: flex;
      height: 100%;
      width: 100%;

      .ds-calendar-event-name {
        margin: auto;
        text-align: center;
      }
    }
  }
</style>
