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

/**
 * Component showing an action icon
 *
 * @module alarm
 *
 * @prop {String} [type] - Type of the action
 * @prop {Function} [method] - Action to execute when user clicks on the action's icon
 * @prop {Boolean} [isDropDown] - Boolean to decide whether to show a dropdown with actions, or actions separately
 */
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
    const eventStyles = { ...this.$constants.EVENT_ENTITY_STYLE };
    const eventTypes = { ...this.$constants.EVENT_ENTITY_TYPES };

    return {
      actionsMap: {
        ack: {
          icon: eventStyles[eventTypes.ack].icon,
          title: this.$t('alarmList.actions.titles.ack'),
        },
        fastAck: {
          icon: eventStyles[eventTypes.fastAck].icon,
          title: this.$t('alarmList.actions.titles.fastAck'),
        },
        ackRemove: {
          icon: eventStyles[eventTypes.ackRemove].icon,
          title: this.$t('alarmList.actions.titles.ackRemove'),
        },
        pbehavior: {
          icon: eventStyles[eventTypes.pbehaviorAdd].icon,
          title: this.$t('alarmList.actions.titles.pbehavior'),
        },
        snooze: {
          icon: eventStyles[eventTypes.snooze].icon,
          title: this.$t('alarmList.actions.titles.snooze'),
        },
        pbehaviorList: {
          icon: eventStyles[eventTypes.pbehaviorList].icon,
          title: this.$t('alarmList.actions.titles.pbehaviorList'),
        },
        declareTicket: {
          icon: eventStyles[eventTypes.declareTicket].icon,
          title: this.$t('alarmList.actions.titles.declareTicket'),
        },
        associateTicket: {
          icon: eventStyles[eventTypes.assocTicket].icon,
          title: this.$t('alarmList.actions.titles.associateTicket'),
        },
        cancel: {
          icon: eventStyles[eventTypes.delete].icon,
          title: this.$t('alarmList.actions.titles.cancel'),
        },
        changeState: {
          icon: eventStyles[eventTypes.changeState].icon,
          title: this.$t('alarmList.actions.titles.changeState'),
        },
        moreInfos: {
          icon: 'more_horiz',
          title: this.$t('alarmList.actions.titles.moreInfos'),
        },
        history: {
          icon: 'history',
          title: this.$t('alarmList.actions.titles.history'),
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
