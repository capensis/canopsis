<template lang="pug">
  span
    template(v-if="isDropDown")
      v-list-tile(@click.stop="method")
        v-list-tile-title
          v-icon.pr-3(left, small) {{ action.icon }}
          span.body-1 {{ $t(action.title) }}
    template(v-else)
      v-btn.mx-1(flat, icon, @click.stop="method", :title="$t(action.title)")
        v-icon {{ action.icon }}
</template>

<script>
const ACTIONS_MAP = {
  ack: {
    icon: 'playlist_add_check',
    title: 'alarmList.actions.ack',
  },
  fastAck: {
    icon: 'check',
    title: 'alarmList.actions.fastAck',
  },
  ackRemove: {
    icon: 'block',
    title: 'alarmList.actions.ackRemove',
  },
  pbehavior: {
    icon: 'pause',
    title: 'alarmList.actions.pbehavior',
  },
  snooze: {
    icon: 'alarm',
    title: 'alarmList.actions.snooze',
  },
  pbehaviorList: {
    icon: 'list',
    title: 'alarmList.actions.pbehaviorList',
  },
  declareTicket: {
    icon: 'local_play',
    title: 'alarmList.actions.declareTicket',
  },
  associateTicket: {
    icon: 'pin_drop',
    title: 'alarmList.actions.associateTicket',
  },
  cancel: {
    icon: 'delete',
    title: 'alarmList.actions.cancel',
  },
  changeState: {
    icon: 'report_problem',
    title: 'alarmList.actions.changeState',
  },
};

export default {
  props: {
    type: {
      validation: value => [Object.keys(ACTIONS_MAP)].includes(value),
      required: true,
    },
    method: {
      type: Function,
      required: true,
    },
    isDropDown: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    action() {
      return ACTIONS_MAP[this.type];
    },
  },
};
</script>
