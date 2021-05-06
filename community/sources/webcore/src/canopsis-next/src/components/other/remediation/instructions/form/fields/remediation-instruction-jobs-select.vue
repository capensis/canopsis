<template lang="pug">
  v-autocomplete(
    v-field="value",
    v-validate="rules",
    :items="jobs",
    :label="label || $t('remediationInstructions.job')",
    :loading="pending",
    :name="name",
    :error-messages="errors.collect(name)",
    item-text="name",
    item-value="_id",
    return-object
  )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import entitiesRemediationJobsMixin from '@/mixins/entities/remediation/jobs';

export default {
  inject: ['$validator'],
  mixins: [entitiesRemediationJobsMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'job',
    },
    label: {
      type: String,
      default: '',
    },
    required: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      jobs: [],
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: jobs } = await this.fetchRemediationJobsListWithoutStore({
        params: {
          limit: MAX_LIMIT,
        },
      });

      this.jobs = jobs;
      this.pending = false;
    },
  },
};
</script>
