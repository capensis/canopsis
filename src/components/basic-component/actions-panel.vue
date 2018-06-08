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
      },
    };
  },
  computed: {
    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: [this.item._id],
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
    createAckEvent() {
      return this.createEvent(EVENT_ENTITY_TYPES.ack, this.item);
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
