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
      let tooltip;

      if (this.alarm.is_manual_instruction_running) {
        tooltip = this.$t('alarmList.tooltips.hasManualInstructionInRunning');
      } else if (this.alarm.is_manual_instruction_waiting_result) {
        tooltip = this.$t('alarmList.tooltips.awaitingInstructionComplete');
      } else if (this.alarm.is_auto_instruction_running) {
        tooltip = this.$t('alarmList.tooltips.hasAutoInstructionInRunning');
      }

      if (tooltip) {
        return {
          tooltip,

          icon: 'assignment',
          class: 'instruction-icon--auto-running',
        };
      }

      if (this.alarm.is_all_auto_instructions_completed) {
        return {
          icon: 'assignment_late',
          class: 'error--text',
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
