<template>

  <div class="ds-day"
       :class="classesDay"
       @click.stop="add"
       @dragstart.prevent>

    <div class="ds-day-header" :class="classesHeader">

      <a class="ds-dom" href
         :class="classesDayOfMonth"
         @click.stop.prevent="viewDay"
         @mousedown.stop>
        {{ dayOfMonth }}
      </a>

      <span class="ds-first-day" v-if="showMonth">
        {{ month }}
      </span>
    </div>

    <template v-for="(event, i) in visibleEvents">
      <ds-calendar-event
        v-bind="{$scopedSlots}"
        v-on="$listeners"
        :key="event.id"
        :calendar-event="event"
        :calendar="calendar"
        :index="i"
      ></ds-calendar-event>
    </template>

    <div v-if="hasPlaceholder">
      <ds-calendar-event-placeholder
        v-bind="{$scopedSlots}"
        v-on="$listeners"
        :day="day"
        :placeholder="placeholder"
        :placeholder-for-create="placeholderForCreate"
        :calendar="calendar"
        :index="visibleEvents.length"
      ></ds-calendar-event-placeholder>
    </div>

  </div>

</template>

<script>
import { DsDay } from 'dayspan-vuetify/src/components';

export default {
  extends: DsDay,
};
</script>

<style lang="scss">
  .ds-day {
    position: relative;

    .ds-dom {
      border-radius: 12px;
      background-color: white;
      display: inline-block;
      position: relative;
      z-index: 1;
    }

    .ds-day-header {
      z-index: 10;
    }
  }
</style>
