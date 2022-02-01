<template lang="pug">
  v-tooltip(top)
    v-icon.instruction-icon(
      :class="iconData.class",
      slot="activator",
      size="16",
      color="black"
    ) {{ iconData.icon }}
    span {{ iconData.tooltip }}
</template>

<script>
export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    iconData() {
      if (this.alarm.is_auto_instruction_running || this.alarm.is_manual_instruction_waiting_result) {
        return {
          icon: 'assignment',
          class: 'instruction-icon--auto-running',
          tooltip: this.alarm.is_manual_instruction_waiting_result
            ? this.$t('alarmList.tooltips.awaitingInstructionComplete')
            : this.$t('alarmList.tooltips.hasAutoInstructionInRunning'),
        };
      }

      if (this.alarm.is_auto_instruction_failed) {
        return {
          icon: 'assignment_late',
          class: 'error--text',
          tooltip: this.$t('alarmList.tooltips.autoInstructionsFailed'),
        };
      }

      if (this.alarm.is_all_auto_instructions_completed) {
        return {
          icon: 'warning',
          tooltip: this.$t('alarmList.tooltips.allAutoInstructionExecuted'),
        };
      }

      return {
        icon: 'assignment',
        tooltip: this.$t('alarmList.tooltips.hasInstruction'),
      };
    },
  },
};
</script>

<style lang="scss" scoped>
.instruction-icon--auto-running {
  border: 2px dashed black;
}
</style>
