<template>
  <c-calendar
    ref="calendar"
    :events="events"
    :loading="pending"
    readonly
    hide-details-menu
    @change:pagination="fetchEvents"
  />
</template>

<script>
import { DATETIME_FORMATS } from '@/constants';

import { getMostReadableTextColor } from '@/helpers/color';
import { getPbehaviorColor, isFullDayEvent } from '@/helpers/entities/pbehavior/form';
import {
  convertDateToDateObjectByTimezone,
  convertDateToString,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';
import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

export default {
  inject: ['$system'],
  mixins: [
    entitiesPbehaviorMixin,
    entitiesPbehaviorTimespansMixin,
  ],
  props: {
    entityId: {
      type: String,
      required: false,
    },
  },
  data() {
    return {
      pending: false,
      events: [],
    };
  },
  mounted() {
    this.fetchEvents();
  },
  methods: {
    convertPbehaviorsCalendarToEvents(pbehaviors) {
      return pbehaviors.map((pbehavior, index) => {
        const start = convertDateToDateObjectByTimezone(pbehavior.from, this.$system.timezone);
        const end = pbehavior.to && convertDateToDateObjectByTimezone(pbehavior.to, this.$system.timezone);

        const isTimed = !isFullDayEvent(start, end);

        const fromString = convertDateToString(pbehavior.from, DATETIME_FORMATS.medium);
        const toString = convertDateToString(pbehavior.to, DATETIME_FORMATS.medium);
        const color = getPbehaviorColor(pbehavior);
        const iconColor = getMostReadableTextColor(color, { level: 'AA', size: 'large' });

        return {
          id: index,
          color,
          iconColor,
          start,
          end,
          icon: pbehavior.type.icon_name,
          name: `${fromString} - ${toString} ${pbehavior.title}`,
          timed: isTimed,
        };
      });
    },

    fetchPbehaviorsCalendar() {
      const params = {
        from: convertDateToTimestampByTimezone(this.$refs.calendar.filled.start, this.$system.timezone),
        to: convertDateToTimestampByTimezone(this.$refs.calendar.filled.end, this.$system.timezone),
      };

      if (this.entityId) {
        params._id = this.entityId;

        return this.fetchEntitiesPbehaviorsCalendarWithoutStore({ params });
      }

      return this.fetchPbehaviorsCalendarWithoutStore({ params });
    },

    async fetchEvents() {
      this.pending = true;

      const pbehaviorsCalendar = await this.fetchPbehaviorsCalendar();

      this.events = this.convertPbehaviorsCalendarToEvents(pbehaviorsCalendar);

      this.pending = false;
    },
  },
};
</script>
