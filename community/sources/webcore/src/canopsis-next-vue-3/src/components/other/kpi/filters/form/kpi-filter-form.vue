<template>
  <c-layout class="gap-3" column>
    <c-name-field v-model="name" autofocus required />
    <kpi-filter-patterns-form :form="patterns" @update:form="patterns = $event" />
  </c-layout>
</template>

<script>
import { useFormModels } from '@/hooks/form'

import KpiFilterPatternsForm from './kpi-filter-patterns-form.vue'

export default {
  components: { KpiFilterPatternsForm },

  setup() {
    const { name, patterns } = useFormModels(['name', 'patterns'])

    return {
      name,
      patterns,
    }
  },
}
</script>
