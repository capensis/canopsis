<template>
  <c-page @refresh="refresh">
    <job-management-tabs ref="tabsElement" />
  </c-page>
</template>

<script>
import { ref } from 'vue';

import JobManagementTabs from '@/components/other/job-management/job-management-tabs.vue';

export default {
  components: {
    JobManagementTabs,
  },

  setup() {
    const tabsElement = ref(null);
    const refresh = () => tabsElement.value?.refresh();

    return {
      tabsElement,
      refresh,
    };
  },
};
</script>
