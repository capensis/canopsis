<template lang="pug">
  div
    div(v-show="$mq === 'laptop'")
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
    div(v-show="$mq === 'mobile' || $mq === 'tablet'")
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
import ModalMixin from '@/mixins/modal/modal';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, ENTITIES_TYPES, ENTITIES_STATUSES, MODALS } from '@/constants';

import ActionsPanelItem from './actions-panel-item.vue';

export default {
  components: { ActionsPanelItem },
  mixins: [ModalMixin, EventActionsMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      actionsMap: {
        ack: {
          icon: 'playlist_add_check',
          title: 'alarmList.actions.ack',
          method: this.showActionModal(MODALS.createAckEvent),
        },
        fastAck: {
          icon: 'check',
          title: 'alarmList.actions.fastAck',
          method: this.createAckEvent,
        },
        ackRemove: {
          icon: 'block',
          title: 'alarmList.actions.ackRemove',
          method: this.showAckRemoveModal,
        },
        pbehavior: {
          icon: 'pause',
          title: 'alarmList.actions.pbehavior',
          method: this.showActionModal(MODALS.createPbehavior),
        },
        snooze: {
          icon: 'alarm',
          title: 'alarmList.actions.snooze',
          method: this.showActionModal(MODALS.createSnoozeEvent),
        },
        pbehaviorList: {
          icon: 'list',
          title: 'alarmList.actions.pbehaviorList',
          method: this.showActionModal(MODALS.pbehaviorList),
        },
        declareTicket: {
          icon: 'local_play',
          title: 'alarmList.actions.declareTicket',
          method: this.showActionModal(MODALS.createDeclareTicketEvent),
        },
        associateTicket: {
          icon: 'pin_drop',
          title: 'alarmList.actions.associateTicket',
          method: this.showActionModal(MODALS.createAssociateTicketEvent),
        },
        cancel: {
          icon: 'delete',
          title: 'alarmList.actions.cancel',
          method: this.showActionModal(MODALS.createCancelEvent),
        },
        changeState: {
          icon: 'report_problem',
          title: 'alarmList.actions.changeState',
          method: this.showActionModal(MODALS.createChangeStateEvent),
        },
      },
    };
  },
  computed: {
    modalConfig() {
      return {
        itemType: ENTITIES_TYPES.alarm,
        itemId: this.item._id,
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
            ],
          };
        }

        return {
          main: [actionsMap.ack, actionsMap.fastAck],
          dropDown: [],
        };
      } else if (this.item.v.status.val === ENTITIES_STATUSES.cancelled) {
        return { // TODO: add restore alarm action
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
    async createAckEvent() {
      await this.createEvent(EVENT_ENTITY_TYPES.ack, this.item);
    },

    showActionModal(name) {
      return () => this.showModal({
        name,
        config: this.modalConfig,
      });
    },

    showAckRemoveModal() {
      this.showModal({
        name: MODALS.createCancelEvent,
        config: {
          ...this.modalConfig,
          title: 'modals.createAckRemove.title',
          eventType: EVENT_ENTITY_TYPES.ackRemove,
        },
      });
    },
  },
};
</script>

<style scoped>
</style>
