<template lang="pug">
  span
    template(v-if="isDropDown")
      v-list-tile(@click.stop="method")
        v-list-tile-title
          v-icon.pr-3(left, small) {{ action.icon }}
          span.body-1 {{ action.title }}
    template(v-else)
      v-btn.mx-1(flat, icon, @click.stop="method", :title="action.title")
        v-icon {{ action.icon }}
</template>

<script>
export default {
  props: {
    type: {
      type: String,
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
  data() {
    return {
      actionsMap: {
        ack: {
          icon: 'playlist_add_check',
          title: this.$t('alarmList.actions.ack'),
        },
        fastAck: {
          icon: 'check',
          title: this.$t('alarmList.actions.fastAck'),
        },
        ackRemove: {
          icon: 'block',
          title: this.$t('alarmList.actions.ackRemove'),
        },
        pbehavior: {
          icon: 'pause',
          title: this.$t('alarmList.actions.pbehavior'),
        },
        snooze: {
          icon: 'alarm',
          title: this.$t('alarmList.actions.snooze'),
        },
        pbehaviorList: {
          icon: 'list',
          title: this.$t('alarmList.actions.pbehaviorList'),
        },
        declareTicket: {
          icon: 'local_play',
          title: this.$t('alarmList.actions.declareTicket'),
        },
        associateTicket: {
          icon: 'pin_drop',
          title: this.$t('alarmList.actions.associateTicket'),
        },
        cancel: {
          icon: 'delete',
          title: this.$t('alarmList.actions.cancel'),
        },
        changeState: {
          icon: 'report_problem',
          title: this.$t('alarmList.actions.changeState'),
        },
      },
    };
  },
  computed: {
    action() {
      return this.actionsMap[this.type] || {};
    },
  },
};
</script>
