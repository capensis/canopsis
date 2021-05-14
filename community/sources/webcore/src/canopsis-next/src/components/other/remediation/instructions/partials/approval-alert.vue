<template lang="pug">
  v-card
    v-card-text
      v-layout(row)
        v-flex.text-right(xs1)
          v-layout(row, justify-end)
            v-icon(color="error") notification_important
        v-flex(xs11)
    v-card-actions
      v-layout.py-1(justify-end)
        v-btn(depressed, flat) Dismiss
        v-btn.primary(type="submit") Approve
</template>

<script>
export default {
  props: {
    approval: {
      type: Object,
      required: true,
    },
  },
};
</script>
