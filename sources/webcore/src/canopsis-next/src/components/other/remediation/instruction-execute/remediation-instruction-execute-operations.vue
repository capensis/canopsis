<template lang="pug">
  v-layout(column)
    remediation-instruction-execute-operation(
      v-for="(operation, index) in operations",
      :key="operation.operation_id",
      :operation="operation",
      :operation-number="getOperationNumber(index + 1)"
    )
</template>

<script>
import { getCharByIndex } from '@/helpers/string';

import RemediationInstructionExecuteOperation from './remediation-instruction-execute-operation.vue';

export default {
  components: { RemediationInstructionExecuteOperation },
  props: {
    operations: {
      type: Array,
      default: () => [],
    },
    stepNumber: {
      type: [Number, String],
      required: true,
    },
  },
  methods: {
    getOperationNumber(index) {
      return `${this.stepNumber}${getCharByIndex(index)}`;
    },
  },
};
</script>
