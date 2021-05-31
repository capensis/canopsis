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
      if (this.alarm.is_auto_instruction_running) {
        return {
          icon: 'assignment',
          class: 'auto-running',
          tooltip: this.$t('alarmList.tooltips.hasAutoInstructionInRunning'),
        };
      } else if (this.alarm.is_all_auto_instructions_completed) {
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
.instruction-icon.auto-running {
  border: 2px dashed black;
}
</style>
