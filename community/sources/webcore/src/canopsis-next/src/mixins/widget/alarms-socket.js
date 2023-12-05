import { isEqual } from 'lodash';
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
      if (!this.liveWatching) {
        return;
      }

      if (!isEqual(mapIds(alarms), mapIds(prevAlarms))) {
        this.leaveAlarmsSocketRoom();
        this.joinToAlarmsSocketRoom(alarms);
      }
    },

    allAlarmDetailsQueries(queries, prevQueries) {
      if (!this.liveWatching || this.editing) {
        return;
      }

      if (!isEqual(mapIds(queries), mapIds(prevQueries))) {
        this.leaveAlarmDetailsSocketRoom();
        this.joinToAlarmDetailsSocketRoom(queries);
      }
    },

    liveWatching: 'toggleSubscription',
    visible: 'toggleSubscription',
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

    toggleSubscription() {
      if (this.visible && this.liveWatching) {
        this.joinToAlarmsSocketRoom(this.alarms);
        this.joinToAlarmDetailsSocketRoom(this.allAlarmDetailsQueries);

        return;
      }

      this.leaveAlarmsSocketRoom();
      this.leaveAlarmDetailsSocketRoom();
    },
  },
};
