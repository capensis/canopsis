<template lang="pug">
  v-container
    v-layout.white(wrap, justify-space-between, align-center)
      v-flex
        v-btn(icon, @click="showSettings")
          v-icon settings
    v-layout.white
      ds-calendar(:events="events")
</template>

<script>
import get from 'lodash/get';
import { createNamespacedHelpers } from 'vuex';

import { SIDE_BARS } from '@/constants';

import sideBarMixin from '@/mixins/side-bar/side-bar';
import widgetQueryMixin from '@/mixins/widget/query';
import { convertAlarmsToCalendarEvents, convertPbehaviorsToCalendarEvents } from '@/helpers/stats/calendar';

import DsCalendar from './day-span/calendar.vue';

const { mapActions: entityMapActions } = createNamespacedHelpers('entity');
const { mapActions: alarmMapActions } = createNamespacedHelpers('alarm');
const { mapActions: pbehaviorMapActions } = createNamespacedHelpers('pbehavior');

export default {
  components: { DsCalendar },
  mixins: [sideBarMixin, widgetQueryMixin],
  props: {
    widget: {
      type: Object,
      required: true,
    },
    rowId: {
      type: String,
      required: true,
    },
  },
  data() {
    return {
      events: [],
    };
  },
  methods: {
    ...entityMapActions({
      fetchContextEntitiesListWithoutStore: 'fetchListWithoutStore',
    }),

    ...alarmMapActions({
      fetchAlarmsListWithoutStore: 'fetchListWithoutStore',
    }),

    ...pbehaviorMapActions({
      fetchPbehaviorsListByEntityIdWithoutStore: 'fetchListByEntityIdWithoutStore',
    }),

    showSettings() {
      this.showSideBar({
        name: SIDE_BARS.statsCalendarSettings,
        config: {
          widget: this.widget,
          rowId: this.rowId,
        },
      });
    },

    async fetchList() {
      const widgetFilter = get(this.widget, 'parameters.mfilter.filter');

      if (widgetFilter) {
        const { entities } = await this.fetchContextEntitiesListWithoutStore({
          params: {
            start: 0,
            limit: 50,
            _filter: widgetFilter,
          },
        });

        const alarmsFilter = {
          $or: [{
            connector_name: {
              $in: entities.map(v => v.name),
            },
          }],
        };

        const { alarms } = await this.fetchAlarmsListWithoutStore({
          params: {
            filter: alarmsFilter,
            skip: 0,
            limit: 15,
          },
        });

        const pbehaviorsCollections = await Promise.all(entities.map(({ _id }) =>
          this.fetchPbehaviorsListByEntityIdWithoutStore({ id: _id })));

        const pbehaviors = [].concat(...pbehaviorsCollections);

        this.events = [
          ...convertAlarmsToCalendarEvents(alarms),
          ...convertPbehaviorsToCalendarEvents(pbehaviors),
        ];
      }
    },
  },
};
</script>
