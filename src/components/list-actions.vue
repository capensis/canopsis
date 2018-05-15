<template lang="pug">
  v-flex
    strong {{ item.pbehaviors.length }}
    v-btn(flat, icon, @click.stop="showAddAckEventModal")
      v-icon playlist_add_check
    v-btn(flat, icon)
      v-icon check
    v-btn(flat, icon, @click.stop="showAddChangeStateEventModal")
      v-icon report_problem
    v-btn(flat, icon, @click.stop="showAddSnoozeEventModal")
      v-icon alarm
    v-btn(flat, icon, @click.stop="showCancelEventModal")
      v-icon delete
    v-btn(flat, icon, @click.stop="showAddPbehaviorModal")
      v-icon pause
    v-btn(flat, icon, @click.stop="showAddPbehaviorModal")
      v-icon local_play
    v-btn(flat, icon, @click.stop="showPbehaviorListModal")
      v-icon list
    v-menu(bottom)
      v-btn(flat, icon, slot="activator")
        v-icon more_vert
      v-list
        v-list-tile
          v-list-tile-action
            v-icon check
          v-list-tile-title Title
        v-list-tile
          v-list-tile-title Something
        v-list-tile
          v-list-tile-title Something
        v-list-tile
          v-list-tile-title Something
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

const { mapActions: modalMapActions } = createNamespacedHelpers('modal');

export default {
  props: {
    item: {
      type: Object,
      required: true,
    },
  },
  methods: {
    ...modalMapActions({
      showModal: 'show',
    }),

    showAddAckEventModal() {
      this.showModal({ name: 'add-ack-event' });
    },
    showCancelEventModal() {
      this.showModal({ name: 'add-cancel-event' });
    },
    showAddChangeStateEventModal() {
      this.showModal({ name: 'add-change-state-event' });
    },
    showAddSnoozeEventModal() {
      this.showModal({ name: 'add-snooze-event' });
    },
    showAddPbehaviorModal() {
      this.showModal({ name: 'add-pbehavior' });
    },
    showPbehaviorListModal() {
      this.showModal({
        name: 'pbehavior-list',
        config: {
          entityKey: 'alarm',
          entityId: this.item._id,
        },
      });
    },
  },
};
</script>
