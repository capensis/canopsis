import { differenceBy } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';

import { mapIds } from '@/helpers/array';

const { mapActions: mapAlarmsActions } = createNamespacedHelpers('alarm');
const { mapGetters: mapAlarmDetailsGetters, mapActions: mapAlarmDetailsActions } = createNamespacedHelpers('alarm/details');

export const widgetAlarmsSocketMixin = {
  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapAlarmDetailsGetters({
      getAlarmDetailsQuery: 'getQuery',
    }),

    alarmsSocketRoom() {
      return `${SOCKET_ROOMS.alarms}/${this.widget._id}`;
    },

    alarmDetailsSocketRoom() {
      return `${SOCKET_ROOMS.alarmDetails}/${this.widget._id}`;
    },

    allAlarmDetailsQueries() {
      return this.getAlarmDetailsQuery(this.widget._id);
    },
  },
  watch: {
    alarms(alarms, prevAlarms) {
      const diff = differenceBy(alarms, prevAlarms, ['_id']);

      if (diff.length) {
        this.leaveAlarmsSocketRoom();
        this.joinToAlarmsSocketRoom(alarms);
      }
    },

    allAlarmDetailsQueries(queries, prevQueries) {
      const diff = differenceBy(queries, prevQueries, ['_id']);

      if (diff.length) {
        this.leaveAlarmDetailsSocketRoom();
        this.joinToAlarmDetailsSocketRoom(queries);
      }
    },
  },
  beforeDestroy() {
    this.leaveAlarmsSocketRoom();
  },
  methods: {
    ...mapAlarmsActions({
      updateAlarmInStore: 'updateItemInStore',
    }),

    ...mapAlarmDetailsActions({
      updateAlarmDetailsInStore: 'updateItemInStore',
    }),

    joinToAlarmsSocketRoom(alarms) {
      this.$socket
        .join(this.alarmsSocketRoom, { ids: mapIds(alarms) })
        .addListener(this.updateAlarmInStore);
    },

    leaveAlarmsSocketRoom() {
      this.$socket
        .leave(this.alarmsSocketRoom)
        .removeListener(this.updateAlarmInStore);
    },

    joinToAlarmDetailsSocketRoom(queries) {
      this.$socket
        .join(this.alarmDetailsSocketRoom, queries)
        .addListener(this.updateAlarmDetailsInStore);
    },

    leaveAlarmDetailsSocketRoom() {
      this.$socket
        .leave(this.alarmDetailsSocketRoom)
        .removeListener(this.updateAlarmDetailsInStore);
    },
  },
};
