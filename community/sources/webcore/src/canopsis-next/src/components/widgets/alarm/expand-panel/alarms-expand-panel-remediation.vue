<template>
  <remediation-instruction-executions-list
    :executions="executions"
    :options.sync="options"
    :total-items="meta.total_count"
    :pending="pending"
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { localQueryMixin } from '@/mixins/query-local/query';

import RemediationInstructionExecutionsList
  from '@/components/other/remediation/instruction-execute/remediation-instruction-executions-list.vue';

const { mapActions } = createNamespacedHelpers('alarm');

export default {
  components: { RemediationInstructionExecutionsList },
  mixins: [localQueryMixin],
  props: {
    alarm: {
      type: Object,
      required: false,
    },
  },
  data() {
    return {
      pending: false,
      meta: {},
      executions: [],
    };
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchAlarmExecutionsWithoutStore: 'fetchExecutionsWithoutStore',
    }),

    async fetchList() {
      try {
        this.pending = true;

        const { data, meta } = await this.fetchAlarmExecutionsWithoutStore({
          id: this.alarm._id,
          params: this.getQuery(),
        });

        this.executions = data;
        this.meta = meta;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
