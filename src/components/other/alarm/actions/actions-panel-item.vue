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
import { EVENT_ENTITY_STYLE, EVENT_ENTITY_TYPES } from '@/constants';

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
    return {
      actionsMap: {
        ack: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ack].icon,
          title: this.$t('alarmList.actions.titles.ack'),
        },
        fastAck: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.fastAck].icon,
          title: this.$t('alarmList.actions.titles.fastAck'),
        },
        ackRemove: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.ackRemove].icon,
          title: this.$t('alarmList.actions.titles.ackRemove'),
        },
        pbehavior: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorAdd].icon,
          title: this.$t('alarmList.actions.titles.pbehavior'),
        },
        snooze: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.snooze].icon,
          title: this.$t('alarmList.actions.titles.snooze'),
        },
        pbehaviorList: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.pbehaviorList].icon,
          title: this.$t('alarmList.actions.titles.pbehaviorList'),
        },
        declareTicket: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.declareTicket].icon,
          title: this.$t('alarmList.actions.titles.declareTicket'),
        },
        associateTicket: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.assocTicket].icon,
          title: this.$t('alarmList.actions.titles.associateTicket'),
        },
        cancel: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.delete].icon,
          title: this.$t('alarmList.actions.titles.cancel'),
        },
        changeState: {
          icon: EVENT_ENTITY_STYLE[EVENT_ENTITY_TYPES.changeState].icon,
          title: this.$t('alarmList.actions.titles.changeState'),
        },
        moreInfos: {
          icon: 'more_horiz',
          title: this.$t('alarmList.actions.titles.moreInfos'),
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
