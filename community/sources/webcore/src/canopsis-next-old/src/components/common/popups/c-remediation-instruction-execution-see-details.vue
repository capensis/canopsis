<template>
  <strong
    class="cursor-pointer text-decoration-underline"
    @click="showExecutionModal"
  >{{ $t('common.seeDetails') }}</strong>
</template>

<script>
import { toRef } from 'vue';

import {
  useShowExecutionModal,
} from '@/components/other/remediation/instruction-execute/hooks/active-remediation-executions-item';

export default {
  props: {
    execution: {
      type: Object,
      required: true,
    },
  },
  setup(props) {
    const { showExecutionModal } = useShowExecutionModal(toRef(props, 'execution'));

    return {
      showExecutionModal,
    };
  },
};
</script>
