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
import actionsPanelMixin from '@/mixins/actions-panel';
import entitiesAlarmMixin from '@/mixins/entities/alarm';
import { EVENT_ENTITY_TYPES, ENTITIES_TYPES, ENTITIES_STATUSES, MODALS } from '@/constants';

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
      actionsMap: {
        ack: {
          type: 'ack',
          method: this.showActionModal(MODALS.createAckEvent),
        },
        fastAck: {
          type: 'fastAck',
          method: this.createAckEvent,
        },
        ackRemove: {
          type: 'ackRemove',
          method: this.showAckRemoveModal,
        },
        pbehavior: {
          type: 'pbehavior',
          method: this.showActionModal(MODALS.createPbehavior),
        },
        snooze: {
          type: 'snooze',
          method: this.showActionModal(MODALS.createSnoozeEvent),
        },
        pbehaviorList: {
          type: 'pbehaviorList',
          method: this.showActionModal(MODALS.pbehaviorList),
        },
        declareTicket: {
          type: 'declareTicket',
          method: this.showActionModal(MODALS.createDeclareTicketEvent),
        },
        associateTicket: {
          type: 'associateTicket',
          method: this.showActionModal(MODALS.createAssociateTicketEvent),
        },
        cancel: {
          type: 'cancel',
          method: this.showActionModal(MODALS.createCancelEvent),
        },
        changeState: {
          type: 'changeState',
          method: this.showActionModal(MODALS.createChangeStateEvent),
        },
        moreInfos: {
          type: 'moreInfos',
          method: this.showMoreInfosModal(),
        },
      },
    };
  },
  computed: {
    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: [this.item._id],
        afterSubmit: () => this.fetchAlarmsListWithPreviousParams({ widgetId: this.widget._id }),
      };
    },
    actions() {
      const { actionsMap } = this;

      if ([ENTITIES_STATUSES.ongoing, ENTITIES_STATUSES.flapping].includes(this.item.v.status.val)) {
        if (this.item.v.ack) {
          return {
            main: [actionsMap.declareTicket, actionsMap.associateTicket, actionsMap.cancel],
            dropDown: [
              actionsMap.ackRemove,
              actionsMap.snooze,
              actionsMap.changeState,
              actionsMap.pbehavior,
              actionsMap.pbehaviorList,
              actionsMap.moreInfos,
            ],
          };
        }

        return {
          main: [actionsMap.ack, actionsMap.fastAck],
          dropDown: [actionsMap.moreInfos],
        };
      } else if (this.item.v.status.val === ENTITIES_STATUSES.cancelled) {
        return {
          main: [actionsMap.pbehaviorList],
          dropDown: [],
        };
      }

      return {
        main: [actionsMap.pbehaviorList],
        dropDown: [],
      };
    },
  },
  methods: {
    createAckEvent() {
      return this.createEvent(EVENT_ENTITY_TYPES.ack, this.item);
    },
  },
};
</script>
