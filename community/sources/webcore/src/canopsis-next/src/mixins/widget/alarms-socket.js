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
      getAlarmDetailsQueries: 'getQueries',
    }),

    alarmsSocketRoom() {
      return `${SOCKET_ROOMS.alarms}/${this.widget._id}`;
    },

    alarmDetailsSocketRoom() {
      return `${SOCKET_ROOMS.alarmDetails}/${this.widget._id}`;
    },

    allAlarmDetailsQueries() {
      return this.getAlarmDetailsQueries(this.widget._id);
    },

    liveWatching() {
      return this.widget.parameters.liveWatching;
    },
  },
  watch: {
    alarms(alarms, prevAlarms) {
      if (!this.liveWatching || this.editing) {
        return;
      }

      const diff = differenceBy(alarms, prevAlarms, '_id');

      if (diff.length) {
        if (prevAlarms.length) {
          this.leaveAlarmsSocketRoom();
        }

        this.joinToAlarmsSocketRoom(alarms);
      }
    },

    allAlarmDetailsQueries(queries, prevQueries) {
      if (!this.liveWatching || this.editing) {
        return;
      }

      const diff = differenceBy(queries, prevQueries, '_id');

      if (diff.length) {
        if (prevQueries.length) {
          this.leaveAlarmDetailsSocketRoom();
        }

        this.joinToAlarmDetailsSocketRoom(queries);
      }
    },

    liveWatching(liveWatching) {
      if (this.editing) {
        return;
      }

      if (liveWatching) {
        this.joinToAlarmsSocketRoom(this.alarms);
        this.joinToAlarmDetailsSocketRoom(this.allAlarmDetailsQueries);

        return;
      }

      this.leaveAlarmsSocketRoom();
      this.leaveAlarmDetailsSocketRoom();
    },
  },
  beforeDestroy() {
    this.leaveAlarmsSocketRoom();
    this.leaveAlarmDetailsSocketRoom();
  },
  methods: {
    ...mapAlarmsActions({
      updateAlarmInStore: 'updateItemInStore',
    }),

    ...mapAlarmDetailsActions({
      updateAlarmDetailsInStore: 'updateItemInStore',
    }),

    joinToAlarmsSocketRoom(alarms) {
      if (this.editing) {
        return;
      }

      this.$socket
        .join(this.alarmsSocketRoom, mapIds(alarms))
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
