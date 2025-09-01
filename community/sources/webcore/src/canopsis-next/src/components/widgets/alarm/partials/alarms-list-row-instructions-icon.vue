<template>
  <c-simple-tooltip :content="iconTooltip" top>
    <template #activator="{ on }">
      <v-icon
        :class="iconClass"
        class="instruction-icon"
        size="22"
        v-on="on"
      >
        {{ iconName }}
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { computed } from 'vue';

import { INSTRUCTION_EXECUTION_ICONS } from '@/constants';

import {
  isInstructionExecutionIconFailed,
  isInstructionExecutionIconInProgress,
  isInstructionExecutionIconSuccess,
  isInstructionExecutionManual,
  hasInstructionWithoutAnyExecution,
} from '@/helpers/entities/remediation/instruction-execution/form';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    alarm: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { t, tc } = useI18n();

    const alarmInstructionExecutionIcon = computed(() => (
      props.alarm.instruction_execution_icon ?? INSTRUCTION_EXECUTION_ICONS.manualAvailable
    ));

    const hasRunningInstruction = computed(() => (
      isInstructionExecutionIconInProgress(alarmInstructionExecutionIcon.value)
    ));

    const someOneInstructionIsFailed = computed(() => (
      isInstructionExecutionIconFailed(alarmInstructionExecutionIcon.value)
    ));

    const someOneInstructionIsSuccessful = computed(() => (
      isInstructionExecutionIconSuccess(alarmInstructionExecutionIcon.value)
    ));

    const isManualInstructionIcon = computed(() => (
      isInstructionExecutionManual(alarmInstructionExecutionIcon.value)
    ));

    const withoutAnyExecution = computed(() => (
      hasInstructionWithoutAnyExecution(alarmInstructionExecutionIcon.value)
    ));

    const iconName = computed(() => {
      if (withoutAnyExecution.value) {
        return '$vuetify.icons.assignment_warning';
      }

      return isManualInstructionIcon.value ? '$vuetify.icons.manual_instruction' : 'assignment';
    });

    const iconClass = computed(() => {
      const classNames = [];

      if (withoutAnyExecution.value) {
        classNames.push('instruction-icon--warning');
      }

      if (hasRunningInstruction.value) {
        classNames.push('blinking', 'instruction-icon--dotted');
      }

      if (someOneInstructionIsFailed.value) {
        classNames.push('instruction-icon--failed');
      }

      if (someOneInstructionIsSuccessful.value) {
        classNames.push('instruction-icon--completed');
      }

      return classNames.join(' ');
    });

    const iconTooltip = computed(() => {
      const {
        running_manual_instructions: runningManualInstructions,
        running_auto_instructions: runningAutoInstructions,
        failed_manual_instructions: failedManualInstructions,
        failed_auto_instructions: failedAutoInstructions,
        successful_manual_instructions: successfulManualInstructions,
        successful_auto_instructions: successfulAutoInstructions,
        assigned_instructions: assignedInstructions,
      } = props.alarm;

      const tooltips = Object.entries({
        runningManualInstructions,
        runningAutoInstructions,
        failedManualInstructions,
        failedAutoInstructions,
        successfulManualInstructions,
        successfulAutoInstructions,
      }).reduce((acc, [key, instructions]) => {
        if (instructions?.length) {
          acc.push(tc(
            `alarm.tooltips.${key}`,
            instructions.length,
            { title: instructions.join(', ') },
          ));
        }

        return acc;
      }, []);

      if (assignedInstructions?.length) {
        tooltips.push(
          withoutAnyExecution.value
            ? t('alarm.tooltips.withoutAnyExecution')
            : tc('alarm.tooltips.hasManualInstruction', assignedInstructions.length),
        );
      }

      return `<span class="pre-wrap">${tooltips.join('\n')}</span>`;
    });

    return {
      iconName,
      iconClass,
      iconTooltip,
    };
  },
};
</script>

<style lang="scss">
.instruction-icon {
  box-sizing: content-box;
  border-width: 1px;
  border-color: transparent;
  border-style: solid;

  .theme--dark &, .theme--light & {
    color: grey;
  }

  &--completed.theme--light.v-icon,
  &--completed.theme--dark.v-icon {
    color: var(--v-primary-base);
  }

  &--failed.theme--light.v-icon,
  &--failed.theme--dark.v-icon {
    color: var(--v-error-base);
  }

  &--dotted {
    border-style: dotted;
    border-color: currentColor;
  }

  &--with-manual-available {
    border-style: dashed;
    border-color: currentColor;
  }

  &--warning svg {
    color: var(--v-warning-base);
  }
}
</style>
