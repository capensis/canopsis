<template>
  <c-tooltip :content="iconTooltip" html top>
    <template #activator="{ props }">
      <v-icon
        :class="iconClass"
        class="instruction-icon"
        size="22"
        v-bind="props"
      >
        {{ iconName }}
      </v-icon>
    </template>
  </c-tooltip>
</template>

<script>
import { computed } from 'vue';

import { INSTRUCTION_EXECUTION_ICONS } from '@/constants';

import {
  hasInstructionWithoutAnyExecution,
  isInstructionExecutionIconFailed,
  isInstructionExecutionIconInProgress,
  isInstructionExecutionIconSuccess,
  isInstructionExecutionManual,
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
    const { t } = useI18n();

    const alarmInstructionExecutionIcon = computed(
      () => props.alarm.instruction_execution_icon ?? INSTRUCTION_EXECUTION_ICONS.manualAvailable,
    );

    const hasRunningInstruction = computed(
      () => isInstructionExecutionIconInProgress(alarmInstructionExecutionIcon.value),
    );

    const someOneInstructionIsFailed = computed(
      () => isInstructionExecutionIconFailed(alarmInstructionExecutionIcon.value),
    );

    const someOneInstructionIsSuccessful = computed(
      () => isInstructionExecutionIconSuccess(alarmInstructionExecutionIcon.value),
    );

    const isManualInstructionIcon = computed(() => isInstructionExecutionManual(alarmInstructionExecutionIcon.value),
    );

    const withoutAnyExecution = computed(() => hasInstructionWithoutAnyExecution(alarmInstructionExecutionIcon.value),
    );

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
          acc.push(t(`alarm.tooltips.${key}`, instructions.length, {
            title: instructions.join(', '),
          }));
        }

        return acc;
      }, []);

      if (withoutAnyExecution.value) {
        tooltips.push(t('alarm.tooltips.withoutAnyExecution'));
      } else if (assignedInstructions?.length) {
        tooltips.push(t('alarm.tooltips.hasManualInstruction', assignedInstructions.length));
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
  border-color: transparent;
  border-style: solid;
  border-width: 1px;
  box-sizing: content-box;

  &.v-icon {
    color: grey;
  }

  &--completed.v-icon {
    color: rgb(var(--v-theme-primary));
  }

  &--failed.v-icon {
    color: rgb(var(--v-theme-error));
  }

  &--dotted {
    border-color: currentColor;
    border-style: dotted;
  }

  &--with-manual-available {
    border-color: currentColor;
    border-style: dashed;
  }

  &--warning svg {
    color: rgb(var(--v-theme-warning));
  }
}
</style>
