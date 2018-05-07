<template lang="pug">
  v-flex
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
    v-btn(flat, icon, @click.stop="showCancelEventModal")
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
export default {
  methods: {
    showAddAckEventModal() {
      this.$store.dispatch('modal/show', {
        componentName: 'add-ack-event',
      });
    },
    showCancelEventModal() {
      this.$store.dispatch('modal/show', {
        componentName: 'add-cancel-event',
      });
    },
    showAddChangeStateEventModal() {
      this.$store.dispatch('modal/show', {
        componentName: 'add-change-state-event',
      });
    },
    showAddSnoozeEventModal() {
      this.$store.dispatch('modal/show', {
        componentName: 'add-snooze-event',
      });
    },
    showAddPbehaviorModal() {
      this.$store.dispatch('modal/show', {
        componentName: 'add-pbehavior',
      });
    },
  },
};
</script>
