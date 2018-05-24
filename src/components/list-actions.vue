<template lang="pug">
  v-flex
    strong {{ item.pbehaviors.length }}
    div(v-show="item.v.ack")
      v-btn(flat, icon, @click.stop="showAckRemoveModal")
        v-icon block
      v-btn(flat, icon, @click.stop="showActionModal('create-declare-ticket-event')")
        v-icon local_play
      v-btn(flat, icon, @click.stop="showActionModal('create-associate-ticket-event')")
        v-icon pin_drop
      v-btn(flat, icon, @click.stop="showActionModal('create-cancel-event')")
        v-icon delete
      v-btn(flat, icon, @click.stop="showActionModal('create-change-state-event')")
        v-icon report_problem
      v-btn(flat, icon, @click.stop="showActionModal('create-snooze-event')")
        v-icon alarm
      v-btn(flat, icon, @click.stop="showActionModal('pbehavior-list')")
        v-icon list
    div(v-show="!item.v.ack")
      v-btn(flat, icon, @click.stop="showActionModal('create-ack-event')")
        v-icon playlist_add_check
      v-btn(flat, icon, @click.stop="createAckEvent")
        v-icon check
      v-btn(flat, icon, @click.stop="showActionModal('create-pbehavior')")
        v-icon pause
      v-btn(flat, icon, @click.stop="showActionModal('create-snooze-event')")
        v-icon alarm
      v-btn(flat, icon, @click.stop="showActionModal('pbehavior-list')")
        v-icon list
</template>

<script>
import ModalMixin from '@/mixins/modal/modal';
import EventActionsMixin from '@/mixins/event-actions';
import { EVENT_TYPES } from '@/config';

/**
 * @prop {Object} item - item of the entity
 */
export default {
  mixins: [ModalMixin, EventActionsMixin],
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  computed: {
    modalConfig() {
      return {
        itemType: 'alarm',
        itemId: this.item._id,
      };
    },
  },
  methods: {
    async createAckEvent() {
      await this.createEvent(EVENT_TYPES.ack, this.item);
    },

    showActionModal(name) {
      this.showModal({
        name,
        config: this.modalConfig,
      });
    },

    showAckRemoveModal() {
      this.showModal({
        name: 'create-cancel-event',
        config: {
          ...this.modalConfig,
          title: 'modals.createAckRemove.title',
          eventType: EVENT_TYPES.ackRemove,
        },
      });
    },
  },
};
</script>
