<template lang="pug">
  span
    template(v-if="isDropDown")
      v-list-tile(@click.stop="method")
        v-list-tile-title
          v-icon.pr-3(left, small) {{ action.icon }}
          span.body-1 {{ action.title }}
    template(v-else)
      v-tooltip(bottom)
        v-btn.mx-1(slot="activator", flat, icon, @click.stop="method")
          v-icon {{ action.icon }}
        span {{ action.title }}
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
          title: this.$t('alarmList.actions.titles.ack'),
        },
        fastAck: {
          icon: 'check',
          title: this.$t('alarmList.actions.titles.fastAck'),
        },
        ackRemove: {
          icon: 'block',
          title: this.$t('alarmList.actions.titles.ackRemove'),
        },
        pbehavior: {
          icon: 'pause',
          title: this.$t('alarmList.actions.titles.pbehavior'),
        },
        snooze: {
          icon: 'alarm',
          title: this.$t('alarmList.actions.titles.snooze'),
        },
        pbehaviorList: {
          icon: 'list',
          title: this.$t('alarmList.actions.titles.pbehaviorList'),
        },
        declareTicket: {
          icon: 'local_play',
          title: this.$t('alarmList.actions.titles.declareTicket'),
        },
        associateTicket: {
          icon: 'pin_drop',
          title: this.$t('alarmList.actions.titles.associateTicket'),
        },
        cancel: {
          icon: 'delete',
          title: this.$t('alarmList.actions.titles.cancel'),
        },
        changeState: {
          icon: 'report_problem',
          title: this.$t('alarmList.actions.titles.changeState'),
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
