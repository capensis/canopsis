<template lang="pug">
  div
    v-card(v-if="!pending")
      v-data-table(
        :items="stats",
        :headers="tableHeaders",
        :pagination.sync="pagination",
        :rows-per-page-items="$config.PAGINATION_PER_PAGE_VALUES",
      )
        template(
          slot="items",
          slot-scope="{ item }",
          xs12,
        )
          td {{ item.entity.name }}
          td
            v-layout(align-center)
              v-chip.px-1(:style="{ backgroundColor: getChipColor(item[query.stat.title].value) }", color="white--text")
                div.body-1.font-weight-bold {{ getChipText(item[query.stat.title].value) }}
              div.caption
                template(v-if="item[query.stat.title].trend >= 0") + {{ item[query.stat.title].trend }}
    v-layout(v-else, justify-center)
      v-progress-circular(
      indeterminate,
      color="primary",
      )
</template>

<script>
import moment from 'moment';

import { PAGINATION_LIMIT } from '@/config';
import { STATS_DISPLAY_MODE, STATS_CRITICITY } from '@/constants';

import { parseStringToDateInterval } from '@/helpers/date-intervals';

import entitiesStatsMixin from '@/mixins/entities/stats';
import widgetQueryMixin from '@/mixins/widget/query';
import entitiesUserPreferenceMixin from '@/mixins/entities/user-preference';

import Ellipsis from '@/components/tables/ellipsis.vue';

export default {
  components: {
    Ellipsis,
  },
  mixins: [
    entitiesStatsMixin,
    widgetQueryMixin,
    entitiesUserPreferenceMixin,
  ],
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      stats: [],
      pagination: {
        sortBy: 'value',
        descending: true,
        rowsPerPage: PAGINATION_LIMIT,
      },
      tableHeaders: [
        {
          text: 'Entity',
          value: 'entity.name',
          sortable: true,
        },
        {
          text: 'Value',
          value: 'value',
          sortable: true,
        },
      ],
    };
  },
  computed: {
    getChipColor() {
      return (value) => {
        const { colors, criticityLevels } = this.widget.parameters.displayMode.parameters;

        if (value < criticityLevels.minor) {
          return colors.ok;
        } else if (value < criticityLevels.major) {
          return colors.minor;
        } else if (value < criticityLevels.critical) {
          return colors.major;
        }

        return colors.critical;
      };
    },
    getChipText() {
      return (value) => {
        const { mode, parameters } = this.widget.parameters.displayMode;
        const { criticityLevels } = parameters;

        if (mode === STATS_DISPLAY_MODE.criticity) {
          if (value < criticityLevels.minor) {
            return STATS_CRITICITY.ok;
          } else if (value < criticityLevels.major) {
            return STATS_CRITICITY.minor;
          } else if (value < criticityLevels.critical) {
            return STATS_CRITICITY.major;
          }
          return STATS_CRITICITY.critical;
        }

        return value;
      };
    },
  },
  methods: {
    // Determine if tstart and tstop are valid Dates or Dynamic Date strings (Ex: 'now')
    dateParse(date, type) {
      if (!moment(date).isValid()) {
        try {
          return parseStringToDateInterval(date, type);
        } catch (err) {
          // TODO: DISPLAY AN ALERT TO THE USER
          console.warn(err);
          return err;
        }
      } else {
        return moment(date);
      }
    },

    async fetchList() {
      this.pending = true;
      const params = {};
      const {
        dateInterval,
        mfilter,
        stat,
        limit,
        sortOrder,
      } = this.getQuery();
      const { periodValue } = dateInterval;
      let { periodUnit, tstart, tstop } = dateInterval;

      tstart = this.dateParse(tstart, 'start');
      tstop = this.dateParse(tstop, 'stop');


      if (periodUnit === 'm') {
        periodUnit = periodUnit.toUpperCase();
        // If period unit is 'month', we need to put the dates at the first day of the month, at 00:00 UTC
        const monthlyRoundedTstart = moment.tz(tstart, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        tstart = monthlyRoundedTstart.add(monthlyRoundedTstart.utcOffset(), 'm');
        const monthlyRoundedTstop = moment.tz(tstop, moment.tz.guess()).startOf('month');
        // Add the difference between the local date, and the UTC one.
        tstop = monthlyRoundedTstop.add(monthlyRoundedTstop.utcOffset(), 'm');
      }

      const stats = {};
      stats[stat.title] = { parameters: stat.parameters, stat: stat.stat.value, trend: true };


      params.duration = `${periodValue}${periodUnit.toLowerCase()}`;
      params.stats = stats;
      params.mfilter = mfilter.filter ? JSON.parse(mfilter.filter) : {};
      params.tstop = tstop.startOf('h').unix();
      params.limit = limit;
      params.sort_column = stat.title;
      params.sort_order = sortOrder.toLowerCase();

      this.stats = await this.fetchStatValuesWithoutStore({
        params,
      });

      this.pending = false;
    },
  },
};
</script>

