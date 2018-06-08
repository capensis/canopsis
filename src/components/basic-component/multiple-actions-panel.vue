<template lang="pug">
  div
    div(v-show="$mq === 'laptop'")
      actions-panel-item(
      v-for="(action, index) in actions",
      v-bind="action",
      :key="`multiple-${index}`"
      )
    div(v-show="$mq === 'mobile' || $mq === 'tablet'")
      v-menu(bottom, left, @click.native.stop)
        v-btn(icon slot="activator")
          v-icon more_vert
        v-list
          actions-panel-item(
          v-for="(action, index) in actions",
          v-bind="action",
          isDropDown,
          :key="`mobile-multiple-${index}`"
          )
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import ModalMixin from '@/mixins/modal/modal';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_ENTITY_TYPES, ENTITIES_TYPES, MODALS } from '@/constants';

import ActionsPanelItem from './actions-panel-item.vue';

const { mapGetters: entitiesMapGetters } = createNamespacedHelpers('entities');

export default {
  components: { ActionsPanelItem },
  mixins: [ModalMixin, EventActionsMixin],
  props: {
    itemsIds: {
      type: Array,
      required: true,
    },
  },
  data() {
    return {
      actions: [
        {
          type: 'pbehavior',
          method: this.showActionModal(MODALS.createPbehavior),
        },
        {
          type: 'fastAck',
          method: this.createAckEvent,
        },
        {
          type: 'ack',
          method: this.showActionModal(MODALS.createAckEvent),
        },
        {
          type: 'ackRemove',
          method: this.showAckRemoveModal,
        },
        // TODO: add restore
      ],
    };
  },
  computed: {
    ...entitiesMapGetters(['getList']),

    items() {
      return this.getList(ENTITIES_TYPES.alarm, this.itemsIds);
    },

    modalConfig() {
      return {
        itemsType: ENTITIES_TYPES.alarm,
        itemsIds: this.itemsIds,
      };
    },
  },
  methods: {
    async createAckEvent() {
      await this.createEvent(EVENT_ENTITY_TYPES.ack, this.items);
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
