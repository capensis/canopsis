<template lang="pug">
  div
    v-select(
      v-if="$mq === 'mobile' || $mq === 'tablet'",
      label="Mass actions",
      :items="actionsMap",
      hide-details
    )
      template(slot="selection" slot-scope="data")
        div Select an action
    v-layout(v-else, align-center, wrap)
      div Selection :
      v-tooltip(v-for="action in actionsMap", :key="action.name", bottom)
        v-btn(slot="activator", small, icon)
          v-icon {{ action.icon }}
        span {{ action.text }}
</template>

<script>
export default {
  name: 'mass-actions',
  data() {
    return {
      actionsMap: [
        {
          name: 'add_pbehaviour',
          text: 'Add a pbehaviour',
          icon: 'pause',
        },
        {
          name: 'fast_ack',
          text: 'Fast ack',
          icon: 'check',
        },
        {
          name: 'ack',
          text: 'ack',
          icon: 'playlist_add_check',
        },
        {
          name: 'cancel_ack',
          text: 'Cancel Ack',
          icon: 'block',
        },
        {
          name: 'restore_alarm',
          text: 'Restore Alarm',
          icon: 'undo',
        },
      ],
    };
  },
};
</script>

