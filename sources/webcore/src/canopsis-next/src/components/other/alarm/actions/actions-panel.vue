<template lang="pug">
  div
      div(v-if="$options.filters.mq($mq, { l: true })")
        v-layout
          actions-panel-item(
          v-for="(action, index) in actions.main",
          v-bind="action",
          :key="`main-${index}`"
          )
          v-menu(v-show="actions.dropDown && Object.keys(actions.dropDown).length", bottom, left, @click.native.stop)
            v-btn(icon, slot="activator")
              v-icon more_vert
            v-list
              actions-panel-item(
              v-for="(action, index) in actions.dropDown",
              v-bind="action",
              isDropDown,
              :key="`drop-down-${index}`"
              )
      div(v-if="$options.filters.mq($mq, { m: true, l: false })")
        v-layout
          v-menu(
          v-show="Object.keys(actions.main).length + Object.keys(actions.dropDown).length > 0",
          bottom,
          left,
          @click.native.stop
          )
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

import { MODALS, ENTITIES_TYPES, ENTITIES_STATUSES, EVENT_ENTITY_TYPES, ALARMLIST_ACTION_PANEL_ACTIONS_TYPES } from '@/constants';

import authMixin from '@/mixins/auth';
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
  mixins: [authMixin, actionsPanelMixin, entitiesAlarmMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
    widget: {
      type: Object,
      required: true,
    },
    isEditingMode: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      allActionsMap: {
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
        variablesHelp: {
          type: 'variablesHelp',
          method: this.showVariablesHelperModal(),
        },
      },
    };
  },
  computed: {
    actionsMap() {
      return pickBy(this.allActionsMap, this.actionsAccessFilterHandler);
    },
    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: [this.item._id],
        afterSubmit: () => this.fetchAlarmsListWithPreviousParams({ widgetId: this.widget._id }),
      };
    },
    actions() {
      if ([ENTITIES_STATUSES.ongoing, ENTITIES_STATUSES.flapping]
        .includes(this.item.v.status.val)) {
        if (this.item.v.ack) {
          return {
            main: pick(
              this.actionsMap,
              [
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.declareTicket,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.declareTicket,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.associateTicket,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.cancel,
                this.isEditingMode ? ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.variablesHelp : null,
              ],
            ),
            dropDown: pick(
              this.actionsMap,
              [
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.ackRemove,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.snooze,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.changeState,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.pbehavior,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.pbehaviorList,
                ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.moreInfos,
              ],
            ),
          };
        }

        return {
          main: pick(
            this.actionsMap,
            [
              ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.ack,
              ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.fastAck,
              this.isEditingMode ? ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.variablesHelp : null,
            ],
          ),
          dropDown: pick(this.actionsMap, [ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.moreInfos]),
        };
      } else if (this.item.v.status.val === ENTITIES_STATUSES.cancelled) {
        return {
          main: pick(this.actionsMap, [ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.pbehaviorList]),
          dropDown: [],
        };
      }

      return {
        main: pick(this.actionsMap, [ALARMLIST_ACTION_PANEL_ACTIONS_TYPES.pbehaviorList]),
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
