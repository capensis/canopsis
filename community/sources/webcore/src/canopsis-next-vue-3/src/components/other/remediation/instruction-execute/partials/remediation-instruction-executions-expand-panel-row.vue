<template>
  <tr>
    <td>
      <v-icon color="grey">
        {{ icon }}
      </v-icon>
    </td>
    <td>
      <span :class="{ 'text-medium-emphasis': greyText }">
        {{ t(`remediation.instructionExecute.stepsTitles.${step.type}.${step.status}`, { name: step.name }) }}
      </span>
      <c-enabled v-if="resultIcon" :value="resultIcon.value" class="ml-1" />
    </td>
    <td>
      {{ step.fail_reason }}
    </td>
    <td>{{ completedAt }}</td>
  </tr>
</template>

<script>
import { computed } from 'vue';

import {
  DATETIME_FORMATS,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES,
  REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES,
} from '@/constants';

import { convertDateToString } from '@/helpers/date/date';

import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    step: {
      type: Object,
      default: () => ({}),
    },
  },

  setup(props) {
    const { t } = useI18n();

    const icon = computed(() => props.step.type === REMEDIATION_INSTRUCTION_EXECUTION_STEP_TYPES.manual
        ? '$vuetify.icons.manual_instruction'
        : 'assignment',
    );
    const greyText = computed(
      () => props.step.status === REMEDIATION_INSTRUCTION_EXECUTION_STEP_STATUSES.skipped,
    );
    const completedAt = computed(() => convertDateToString(props.step.completed_at, DATETIME_FORMATS.long, '-'),
    );
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
      completedAt,
      greyText,
      icon,
      resultIcon,
      t,
    };
  },
};
</script>
