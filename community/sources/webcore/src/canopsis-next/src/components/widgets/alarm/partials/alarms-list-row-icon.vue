<template lang="pug">
  v-tooltip(top)
    template(#activator="{ on }")
      v-icon.instruction-icon(
        v-on="on",
        :class="iconClass",
        size="24",
        color="black"
      ) {{ iconName }}
    span.pre-wrap(v-html="iconTooltip")
</template>

<script>
import {
  isInstructionExecutionExecutedAndOtherAvailable,
  isInstructionExecutionIconFailed,
  isInstructionExecutionIconInProgress,
  isInstructionExecutionIconSuccess,
  isInstructionExecutionManual,
} from '@/helpers/forms/remediation-instruction-execution';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  computed: {
    alarmInstructionExecutionIcon() {
      return this.alarm.instruction_execution_icon;
    },

    hasRunningInstruction() {
      return isInstructionExecutionIconInProgress(this.alarmInstructionExecutionIcon);
    },

    someOneInstructionIsFailed() {
      return isInstructionExecutionIconFailed(this.alarmInstructionExecutionIcon);
    },

    someOneInstructionIsSuccessful() {
      return isInstructionExecutionIconSuccess(this.alarmInstructionExecutionIcon);
    },

    instructionExecutedAndOtherAvailable() {
      return isInstructionExecutionExecutedAndOtherAvailable(this.alarmInstructionExecutionIcon);
    },

    isManualInstructionIcon() {
      return isInstructionExecutionManual(this.alarmInstructionExecutionIcon);
    },

    iconName() {
      if (this.isManualInstructionIcon) {
        return '$vuetify.icons.manual_instruction';
      }

      return 'assignment';
    },

    iconClass() {
      const classNames = [];

      if (this.hasRunningInstruction) {
        classNames.push('blinking', 'instruction-icon--dotted');
      }

      if (this.someOneInstructionIsFailed) {
        classNames.push('error--text');
      }

      if (this.someOneInstructionIsSuccessful) {
        classNames.push('primary--text');
      }

      if (this.instructionExecutedAndOtherAvailable) {
        classNames.push('instruction-icon--dashed');
      }

      return classNames.join(' ');
    },

    iconTooltip() {
      const {
        running_manual_instructions: runningManualInstructions,
        running_auto_instructions: runningAutoInstructions,
        failed_manual_instructions: failedManualInstructions,
        failed_auto_instructions: failedAutoInstructions,
        successful_manual_instructions: successfulManualInstructions,
        successful_auto_instructions: successfulAutoInstructions,
        assigned_instructions: assignedInstructions,
      } = this.alarm;

      const tooltips = Object.entries({
        runningManualInstructions,
        runningAutoInstructions,
        failedManualInstructions,
        failedAutoInstructions,
        successfulManualInstructions,
        successfulAutoInstructions,
      }).reduce((acc, [key, instructions]) => {
        if (instructions?.length) {
          acc.push(this.$t(`alarmList.tooltips.${key}`, { title: instructions.join(', ') }));
        }

        return acc;
      }, []);

      if (assignedInstructions?.length) {
        tooltips.push(this.$t('alarmList.tooltips.hasManualInstruction'));
      }

      return tooltips.join('\n');
    },
  },
};
</script>

<style lang="scss" scoped>
.instruction-icon {
  box-sizing: content-box;
  border-width: 1px;
  border-color: transparent;
  border-style: solid;

  &--dashed {
    border-style: dashed;
    border-color: currentColor;
  }

  &--dotted {
    border-style: dotted;
    border-color: currentColor;
  }
}
</style>
