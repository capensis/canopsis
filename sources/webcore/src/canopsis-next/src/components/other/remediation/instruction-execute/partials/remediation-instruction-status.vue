<template lang="pug">
  v-layout.my-2.black--text(:class="{ 'error--text': isFailedStatus }", row, align-center)
    v-icon.mr-1(:color="iconColor") alarm
    span {{ label }}
</template>

<script>
export default {
  props: {
    failedAt: {
      type: Number,
      required: false,
    },
    completedAt: {
      type: Number,
      required: false,
    },
    startedAt: {
      type: Number,
      required: false,
    },
    timeToComplete: {
      type: Object,
      required: true,
    },
  },
  computed: {
    isFailedStatus() {
      return !!this.failedAt;
    },

    iconColor() {
      return this.isFailedStatus ? 'error' : 'black';
    },

    completedLabel() {
      return this.$t('remediationInstructionExecute.completedAt', {
        time: this.convertDate(this.completedAt),
      });
    },

    startedLabel() {
      return this.$t('remediationInstructionExecute.startedAt', {
        time: this.convertDate(this.startedAt),
      });
    },

    failedLabel() {
      return this.$t('remediationInstructionExecute.failedAt', {
        time: this.convertDate(this.failedAt),
      });
    },

    timeToCompleteLabel() {
      const duration = this.$options.filters.duration(this.timeToComplete.seconds);

      return this.$t('remediationInstructionExecute.timeToComplete', { duration });
    },

    label() {
      if (this.completedAt) {
        return this.completedLabel;
      }

      if (this.failedAt) {
        return this.failedLabel;
      }

      if (this.startedAt) {
        return this.startedLabel;
      }

      return this.timeToCompleteLabel;
    },
  },
  methods: {
    convertDate(date) {
      return this.$options.filters.date(date, 'dateTimePicker', true);
    },
  },
};
</script>
