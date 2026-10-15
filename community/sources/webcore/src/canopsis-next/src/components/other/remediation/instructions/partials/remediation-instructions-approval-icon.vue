<template>
  <c-tooltip :disabled="isApproved" location="top">
    <template #activator="{ props: activatorProps }">
      <span v-bind="activatorProps">
        <c-icon :color="iconData.color">
          {{ iconData.name }}
        </c-icon>
      </span>
    </template>
    {{ iconData.tooltip }}
  </c-tooltip>
</template>

<script>
import {
  isInstructionApproved,
  isInstructionDismissed,
} from '@/helpers/entities/remediation/instruction/form';

export default {
  props: {
    instruction: {
      type: Object,
      required: true,
    },
  },

  computed: {
    isDismissed() {
      return isInstructionDismissed(this.instruction);
    },

    isApproved() {
      return isInstructionApproved(this.instruction);
    },

    iconData() {
      if (this.isApproved) {
        return {
          color: 'primary',
          name: 'check_circle',
        };
      }

      return {
        color: 'black',
        name: 'query_builder',
        tooltip: this.isDismissed
          ? this.$t('remediation.instruction.approvalDismissed')
          : this.$t('remediation.instruction.approvalPending'),
      };
    },
  },
};
</script>
