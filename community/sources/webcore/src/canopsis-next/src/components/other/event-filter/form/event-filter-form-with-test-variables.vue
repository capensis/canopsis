<template>
  <v-tabs fixed-tabs>
    <v-tab>{{ $t('common.general') }}</v-tab>
    <v-tab>{{ $tc('templateTesting.testVariables') }}</v-tab>

    <v-tab-item eager>
      <event-filter-form
        v-field="form"
        v-bind="$attrs"
      />
    </v-tab-item>

    <v-tab-item>
      <template-testing-test-variables v-field="form" />
    </v-tab-item>
  </v-tabs>
</template>

<script>
import { ref, onMounted } from 'vue';

import { templateVarsToVariables } from '@/helpers/variables';

import { usePendingHandler } from '@/hooks/query/pending';
import { useTemplateVars } from '@/hooks/store/modules/template-vars';

import TemplateTestingTestVariables from '@/components/other/template-testing/template-testing-test-variables.vue';

import EventFilterForm from './event-filter-form.vue';

export default {
  components: {
    TemplateTestingTestVariables,
    EventFilterForm,
  },
  inheritAttrs: false,
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
  },
  setup() {
    const templateVars = ref([]);
    const { fetchEventFiltersTemplateVarsWithoutStore } = useTemplateVars();

    const {
      pending,
      handler: fetchList,
    } = usePendingHandler(async () => {
      const data = await fetchEventFiltersTemplateVarsWithoutStore();

      templateVars.value = templateVarsToVariables(data);
    });

    onMounted(fetchList);

    return {
      pending,
    };
  },
};
</script>
