<template>
  <tr>
    <td>
      <v-icon color="grey darken-2">
        {{ icon }}
      </v-icon>
    </td>
    <td>
      <span :class="{ 'grey--text': greyText }">
        {{ $t(`remediation.instructionExecute.stepsTitles.${step.type}.${step.status}`, { name: step.name }) }}
      </span>
      <c-enabled v-if="resultIcon" :value="resultIcon.value" class="ml-1" />
    </td>
    <td>
      {{ step.fail_reason }}
    </td>
    <td>{{ step.completed_at | date('long', '-') }}</td>
  </tr>
</template>

<script>
import { computed } from 'vue';

import {
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES,
} from '@/constants';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },
  setup(props) {
    const icon = computed(() => (
      props.step.type === REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES.manual
        ? '$vuetify.icons.manual_instruction'
        : 'assignment'
    ));
    const greyText = computed(() => props.step.status === REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped);
    const resultIcon = computed(() => {
      switch (props.step.status) {
        case REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.completed:
          return { value: true };

        case REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.failed:
          return { value: false };

        default:
          return null;
      }
    });

    return {
      icon,
      greyText,
      resultIcon,
    };
  },
};
</script>
