<template>
  <c-select
    :model-value="modelValue"
    :items="configurations"
    :label="$t('remediation.job.configuration')"
    :loading="pending"
    name="configuration"
    item-title="name"
    item-value="_id"
    return-object
    required
    @update:model-value="$emit('update:modelValue', $event)"
  />
</template>

<script>
import { onMounted, ref } from 'vue';

import { MAX_LIMIT } from '@/constants';
import { useRemediationConfiguration } from '@/hooks/store/modules/remediation-configuration';

export default {
  props: {
    modelValue: {
      type: [Object, String],
      required: false,
    },
  },

  emits: ['update:modelValue'],

  setup() {
    const { fetchRemediationConfigurationsListWithoutStore } = useRemediationConfiguration();

    const pending = ref(false);
    const configurations = ref([]);

    const fetchConfigurations = async () => {
      try {
        pending.value = true;

        const { data } = await fetchRemediationConfigurationsListWithoutStore({
          params: { limit: MAX_LIMIT },
        });

        configurations.value = data;
      } finally {
        pending.value = false;
      }
    };

    onMounted(fetchConfigurations);

    return {
      pending,
      configurations,
    };
  },
};
</script>
