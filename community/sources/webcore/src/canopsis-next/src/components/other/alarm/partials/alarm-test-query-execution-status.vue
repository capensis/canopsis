<template>
  <span>
    <v-progress-circular
      v-if="running"
      :color="color"
      width="3"
      size="22"
      indeterminate
    />
    <c-help-icon
      v-else
      :icon="statusIcon.name"
      :color="statusIcon.color"
      :text="statusIcon.text"
      top
    />
  </span>
</template>

<script>
export default {
  props: {
    running: {
      type: Boolean,
      default: false,
    },
    success: {
      type: Boolean,
      default: false,
    },
    failReason: {
      type: String,
      default: '',
    },
    color: {
      type: String,
      default: 'primary',
    },
  },
  computed: {
    statusIcon() {
      if (this.success) {
        return {
          name: 'check_circle',
          color: this.color,
          text: this.$t('declareTicket.webhookComplete'),
        };
      }

      return {
        name: 'error',
        color: 'error',
        text: this.failReason,
      };
    },
  },
};
</script>
