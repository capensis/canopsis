import { isEqual } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { SOCKET_ROOMS } from '@/config';

import { mapIds } from '@/helpers/array';
import { convertAlarmDetailsQueryToRequest } from '@/helpers/entities/alarm/query';

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

    alarmDetailsSocketPayload() {
      return this.allAlarmDetailsQueries.map(convertAlarmDetailsQueryToRequest);
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

    alarmDetailsSocketPayload(payload, prevPayload) {
      if (!this.liveWatching || this.editing) {
        return;
      }

      if (isEqual(payload, prevPayload)) {
        return;
      }

      if (!payload.length) {
        this.leaveAlarmDetailsSocketRoom();

        return;
      }

      this.leaveAlarmDetailsSocketRoom();
      this.joinToAlarmDetailsSocketRoom(payload);
    },

    liveWatching: 'toggleSubscription',
    visible: 'toggleSubscription',
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
      this.$socket
        .join(this.alarmsSocketRoom, mapIds(alarms))
        .addListener(this.updateAlarmInStore);
    },

    leaveAlarmsSocketRoom() {
      this.$socket
        .leave(this.alarmsSocketRoom)
        .removeListener(this.updateAlarmInStore);
    },

    joinToAlarmDetailsSocketRoom(payload) {
      if (!payload?.length) {
        return;
      }

      this.$socket
        .join(this.alarmDetailsSocketRoom, payload)
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

        if (this.alarmDetailsSocketPayload.length) {
          this.joinToAlarmDetailsSocketRoom(this.alarmDetailsSocketPayload);
        } else {
          this.leaveAlarmDetailsSocketRoom();
        }

        return;
      }

      this.leaveAlarmsSocketRoom();
      this.leaveAlarmDetailsSocketRoom();
    },
  },
};
