<template lang="pug">
  v-layout(row)
    v-select(
      v-field="value",
      v-validate="'required'",
      :items="configurations",
      :label="$t('modals.createRemediationJob.fields.configuration')",
      :error-messages="errors.collect('configuration')",
      :loading="pending",
      name="configuration",
      return-object,
      item-text="name",
      item-value="_id"
    )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import entitiesRemediationConfigurationsMixin from '@/mixins/entities/remediation/configurations';

export default {
  inject: ['$validator'],
  mixins: [entitiesRemediationConfigurationsMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String],
      required: false,
    },
  },
  data() {
    return {
      pending: false,
      configurations: [],
    };
  },
  mounted() {
    this.fetchConfigurations();
  },
  methods: {
    async fetchConfigurations() {
      this.pending = true;

      const { data: configurations } = await this.fetchRemediationConfigurationsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.configurations = configurations;
      this.pending = false;
    },
  },
};
</script>
