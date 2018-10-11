<template>
  <v-menu
    class="ds-calendar-event-menu"
    :class="menuClass"
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
          <v-layout row>
            <v-flex class="ds-ev-description-wrapper" xs10>
              <span class="ds-ev-description">{{ details.description }}</span>
            </v-flex>
            <v-flex xs2>
              <strong class="ds-ev-title">{{ details.title }}</strong>
            </v-flex>
          </v-layout>
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

    menuClass() {
      return this.calendarEvent.data.meta.type;
    },
  },
  methods: {
    editCheck() {

    },
  },
};
</script>

<style lang="scss">
  .single {
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
  }

  .multiple {
    .ds-ev-description-wrapper {
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }
  }
</style>
