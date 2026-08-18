import { isEqual } from 'lodash';
import { mapActions, mapState } from 'pinia';

import { SOCKET_ROOMS } from '@/config';

import { useAlarmStore } from '@/stores/modules/alarm';

import { mapIds } from '@/helpers/array';
import { convertAlarmDetailsQueryToRequest } from '@/helpers/entities/alarm/query';

import { activeViewMixin } from '@/mixins/active-view';

export const widgetAlarmsSocketMixin = {
  mixins: [activeViewMixin],

  props: {
    widget: {
      type: Object,
      required: true,
    },
  },
  computed: {
    ...mapState(useAlarmStore, {
      allAlarmDetailsQueries(store) {
        return store.getQueriesByWidgetId(this.widget._id);
      },
    }),

    alarmsSocketRoom() {
      return `${SOCKET_ROOMS.alarms}/${this.widget._id}`;
    },

    alarmDetailsSocketRoom() {
      return `${SOCKET_ROOMS.alarmDetails}/${this.widget._id}`;
    },

    alarmDetailsSocketPayload() {
      return this.allAlarmDetailsQueries.filter(Boolean).map(convertAlarmDetailsQueryToRequest);
    },

    liveWatching() {
      return this.widget.parameters.liveWatching;
    },

    editing() {
      return this.activeViewEditing;
    },
  },
  watch: {
    alarms(alarms, prevAlarms) {
      if (!this.liveWatching || !this.visible) {
        return;
      }

      if (!isEqual(mapIds(alarms), mapIds(prevAlarms))) {
        this.leaveAlarmsSocketRoom();
        this.joinToAlarmsSocketRoom(alarms);
      }
    },

    alarmDetailsSocketPayload(payload, prevPayload) {
      if (!this.liveWatching || this.editing || !this.visible) {
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
    editing: 'toggleSubscription',
  },
  beforeUnmount() {
    this.leaveAlarmsSocketRoom();
    this.leaveAlarmDetailsSocketRoom();
  },
  methods: {
    ...mapActions(useAlarmStore, {
      updateAlarmInStore: 'updateItemInStore',
      updateAlarmDetailsInStore: 'updateDetailsItemInStore',
    }),

    handleAlarmDetailsSocketUpdate(alarmDetails) {
      this.updateAlarmDetailsInStore({
        widgetId: this.widget._id,
        details: alarmDetails,
      });
    },

    joinToAlarmsSocketRoom(alarms) {
      this.$socket.join(this.alarmsSocketRoom, mapIds(alarms)).addListener(this.updateAlarmInStore);
    },

    leaveAlarmsSocketRoom() {
      this.$socket.leave(this.alarmsSocketRoom).removeListener(this.updateAlarmInStore);
    },

    joinToAlarmDetailsSocketRoom(payload) {
      if (!payload?.length) {
        return;
      }

      this.$socket
        .join(this.alarmDetailsSocketRoom, payload)
        .addListener(this.handleAlarmDetailsSocketUpdate);
    },

    leaveAlarmDetailsSocketRoom() {
      this.$socket
        .leave(this.alarmDetailsSocketRoom)
        .removeListener(this.handleAlarmDetailsSocketUpdate);
    },

    toggleSubscription() {
      this.leaveAlarmsSocketRoom();
      this.leaveAlarmDetailsSocketRoom();

      if (!this.visible || !this.liveWatching) {
        return;
      }

      this.joinToAlarmsSocketRoom(this.alarms);

      if (this.alarmDetailsSocketPayload.length) {
        this.joinToAlarmDetailsSocketRoom(this.alarmDetailsSocketPayload);
      }
    },
  },
};
