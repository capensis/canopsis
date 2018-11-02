<template lang="pug">
  div
      div(v-show="$options.filters.mq($mq, { l: true })")
        v-layout
          actions-panel-item(
          v-for="(action, index) in actions.main",
          v-bind="action",
          :key="`main-${index}`"
          )
          v-menu(v-show="actions.dropDown && actions.dropDown.length", bottom, left, @click.native.stop)
            v-btn(icon, slot="activator")
              v-icon more_vert
            v-list
              actions-panel-item(
              v-for="(action, index) in actions.dropDown",
              v-bind="action",
              isDropDown,
              :key="`drop-down-${index}`"
              )
      div(v-show="$options.filters.mq($mq, { m: true, l: false })")
        v-layout
          v-menu(bottom, left, @click.native.stop)
            v-btn(icon slot="activator")
              v-icon more_vert
            v-list
              actions-panel-item(
              v-for="(action, index) in actions.main",
              v-bind="action",
              isDropDown,
              :key="`mobile-main-${index}`"
              )
              actions-panel-item(
              v-for="(action, index) in actions.dropDown",
              v-bind="action",
              isDropDown,
              :key="`mobile-drop-down-${index}`"
              )
</template>

<script>
import pick from 'lodash/pick';
import pickBy from 'lodash/pickBy';

import actionsPanelMixin from '@/mixins/actions-panel';
import entitiesAlarmMixin from '@/mixins/entities/alarm';

import ActionsPanelItem from './actions-panel-item.vue';

/**
 * Component to regroup actions (actions-panel-item) for each alarm on the alarms list
 *
 * @module alarm
 *
 * @prop {Object} [item] - Object representing an alarm
 */
export default {
  components: { ActionsPanelItem },
  mixins: [actionsPanelMixin, entitiesAlarmMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      allActionsMap: {
        ack: {
          type: 'ack',
          right: 'listalarm_ack',
          method: this.showActionModal(this.$constants.MODALS.createAckEvent),
        },
        fastAck: {
          type: 'fastAck',
          right: 'listalarm_fastAck',
          method: this.createAckEvent,
        },
        ackRemove: {
          type: 'ackRemove',
          right: 'listalarm_cancelAck',
          method: this.showAckRemoveModal,
        },
        pbehavior: {
          type: 'pbehavior',
          right: 'listalarm_pbehavior',
          method: this.showActionModal(this.$constants.MODALS.createPbehavior),
        },
        snooze: {
          type: 'snooze',
          right: 'listalarm_snoozeAlarm',
          method: this.showActionModal(this.$constants.MODALS.createSnoozeEvent),
        },
        pbehaviorList: {
          type: 'pbehaviorList',
          right: 'listalarm_listPbehavior',
          method: this.showActionModal(this.$constants.MODALS.pbehaviorList),
        },
        declareTicket: {
          type: 'declareTicket',
          right: 'listalarm_declareanIncident',
          method: this.showActionModal(this.$constants.MODALS.createDeclareTicketEvent),
        },
        associateTicket: {
          type: 'associateTicket',
          right: 'listalarm_assignTicketNumber',
          method: this.showActionModal(this.$constants.MODALS.createAssociateTicketEvent),
        },
        cancel: {
          type: 'cancel',
          right: 'listalarm_removeAlarm',
          method: this.showActionModal(this.$constants.MODALS.createCancelEvent),
        },
        changeState: {
          type: 'changeState',
          right: 'listalarm_changeState',
          method: this.showActionModal(this.$constants.MODALS.createChangeStateEvent),
        },
        moreInfos: {
          type: 'moreInfos',
          method: this.showMoreInfosModal(),
        },
      },
    };
  },
  computed: {
    actionsMap() {
      return pickBy(this.allActionsMap, value => this.hasAccess(value.right));
    },
    modalConfig() {
      return {
        itemsType: this.$constants.ENTITIES_TYPES.alarm,
        itemsIds: [this.item._id],
        afterSubmit: () => this.fetchAlarmsListWithPreviousParams({ widgetId: this.widget._id }),
      };
    },
    actions() {
      const { actionsMap } = this;

      if ([this.$constants.ENTITIES_STATUSES.ongoing, this.$constants.ENTITIES_STATUSES.flapping]
        .includes(this.item.v.status.val)) {
        if (this.item.v.ack) {
          return {
            main: pick(actionsMap, ['declareTicket', 'associateTicket', 'cancel']),
            dropDown: pick(actionsMap, [
              'ackRemove',
              'snooze',
              'changeState',
              'pbehavior',
              'pbehaviorList',
              'moreInfos',
            ]),
          };
        }

        return {
          main: pick(actionsMap, ['ack', 'fastAck']),
          dropDown: pick(actionsMap, ['moreInfos']),
        };
      } else if (this.item.v.status.val === this.$constants.ENTITIES_STATUSES.cancelled) {
        return {
          main: pick(actionsMap, ['pbehaviorList']),
          dropDown: [],
        };
      }

      return {
        main: pick(actionsMap, ['pbehaviorList']),
        dropDown: [],
      };
    },
  },
  methods: {
    createAckEvent() {
      return this.createEvent(this.$constants.EVENT_ENTITY_TYPES.ack, this.item);
    },
  },
};
</script>
