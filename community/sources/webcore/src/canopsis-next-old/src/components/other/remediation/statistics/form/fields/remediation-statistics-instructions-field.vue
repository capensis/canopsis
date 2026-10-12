<template>
  <v-select
    v-field="value"
    :items="items"
    :label="$t('common.instructions')"
    :loading="pending"
    item-text="name"
    item-value="_id"
    hide-details
  />
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT, REMEDIATION_INSTRUCTION_TYPES } from '@/constants';

const { mapActions } = createNamespacedHelpers('remediationInstruction');

export default {
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Number],
      required: true,
    },
  },
  data() {
    return {
      pending: false,
      instructions: [],
    };
  },
  computed: {
    items() {
      if (!this.instructions.length) {
        return [];
      }

      return [
        {
          _id: '',
          name: this.$t('remediation.statistic.allInstructions'),
        },
        {
          _id: REMEDIATION_INSTRUCTION_TYPES.manual,
          name: this.$t('remediation.statistic.manualInstructions'),
        },
        {
          _id: REMEDIATION_INSTRUCTION_TYPES.auto,
          name: this.$t('remediation.statistic.autoInstructions'),
        },

        {
          divider: true,
        },

        ...this.instructions,
      ];
    },
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    ...mapActions({
      fetchInstructionsListWithoutStore: 'fetchListWithoutStore',
    }),

    async fetchList() {
      try {
        this.pending = true;

        const { data } = await this.fetchInstructionsListWithoutStore({
          params: { limit: MAX_LIMIT },
        });

        this.instructions = data;
      } catch (err) {
        console.error(err);
      } finally {
        this.pending = false;
      }
    },
  },
};
</script>
