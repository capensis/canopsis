<template>
  <c-simple-tooltip :content="tooltip">
    <template #activator="{ on }">
      <v-icon :color="color" v-on="on">
        {{ icon }}
      </v-icon>
    </template>
  </c-simple-tooltip>
</template>

<script>
import { computed } from 'vue';

import { COLORS } from '@/config';
import {
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES,
  REMEDIATION_INSTRUCTION_TYPES,
} from '@/constants';

import { isInstructionTypeAnySimpleManual } from '@/helpers/entities/remediation/instruction/form';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    type: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_TYPES.manual,
    },
    status: {
      type: Number,
      default: REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.running,
    },
    name: {
      type: String,
      default: '',
    },
  },
  setup(props) {
    const { t } = useI18n();

    const icon = computed(() => {
      if (props.status === REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped) {
        return '$vuetify.icons.assignment_one';
      }

      return isInstructionTypeAnySimpleManual(props.type) ? '$vuetify.icons.manual_instruction' : 'assignment';
    });

    const color = computed(() => ({
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.running]: COLORS.remediation.executionStatus.running,
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.completed]: COLORS.remediation.executionStatus.completed,
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.failed]: COLORS.remediation.executionStatus.failed,
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.aborted]: COLORS.remediation.executionStatus.aborted,
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped]: COLORS.remediation.executionStatus.skipped,
      [REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.waiting]: COLORS.remediation.executionStatus.running,
    }[props.status]));

    const tooltip = computed(() => {
      const messageKeySuffix = `${REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES.manual}.${props.status}`;

      return t(`remediation.instructionExecute.stepsTitles.${messageKeySuffix}`, { name: props.name });
    });

    return {
      icon,
      color,
      tooltip,
    };
  },
};
</script>
