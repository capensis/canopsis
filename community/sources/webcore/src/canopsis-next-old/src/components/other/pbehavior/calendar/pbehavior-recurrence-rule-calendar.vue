<template>
  <c-calendar
    ref="calendar"
    :events="events"
    :loading="pending"
    readonly
    hide-details-menu
    @change:pagination="fetchTimespans"
  />
</template>

<script>
import { COLORS } from '@/config';
import { DATETIME_FORMATS } from '@/constants';

import { pbehaviorToTimespanRequest } from '@/helpers/entities/pbehavior/timespans/form';
import {
  convertDateToDateObjectByTimezone,
  convertDateToString,
  convertDateToTimestampByTimezone,
} from '@/helpers/date/date';
import { isFullDayEvent } from '@/helpers/entities/pbehavior/form';

import { entitiesInfoMixin } from '@/mixins/entities/info';
import { entitiesPbehaviorTimespansMixin } from '@/mixins/entities/pbehavior/timespans';

export default {
  inject: ['$system'],
  mixins: [
    entitiesInfoMixin,
    entitiesPbehaviorTimespansMixin,
  ],
  props: {
    pbehavior: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      timespans: [],
    };
  },
  computed: {
    isTimed() {
      return !isFullDayEvent(
        convertDateToDateObjectByTimezone(this.pbehavior.tstart, this.$system.timezone),
        this.pbehavior.tstop && convertDateToDateObjectByTimezone(this.pbehavior.tstop, this.$system.timezone),
      );
    },

    events() {
      return this.timespans.map((timespan, index) => {
        const start = convertDateToDateObjectByTimezone(timespan.from, this.$system.timezone);
        const end = convertDateToDateObjectByTimezone(timespan.to, this.$system.timezone);

        const fromString = convertDateToString(timespan.from, DATETIME_FORMATS.medium);
        const toString = convertDateToString(timespan.to, DATETIME_FORMATS.medium);

        return {
          id: index,
          start,
          end,
          timed: this.isTimed,
          name: `${fromString} - ${toString}`,
          color: COLORS.secondary,
        };
      });
    },
  },
  mounted() {
    this.fetchTimespans();
  },
  methods: {
    async fetchTimespans() {
      this.pending = true;

      const { start, end } = this.$refs.calendar.filled;

      const from = convertDateToTimestampByTimezone(start, this.$system.timezone);
      const to = convertDateToTimestampByTimezone(end, this.$system.timezone);

      const timespan = pbehaviorToTimespanRequest({
        pbehavior: this.pbehavior,
        from,
        to,
      });

      this.timespans = await this.fetchTimespansListWithoutStore({ data: timespan });

      this.pending = false;
    },
  },
};
</script>
