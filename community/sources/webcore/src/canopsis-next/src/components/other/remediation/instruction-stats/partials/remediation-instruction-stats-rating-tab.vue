<template>
  <c-advanced-data-table
    :items="remediationInstructionStatsComments"
    :headers="headers"
    :loading="pending"
    :options.sync="options"
    :total-items="totalItems"
    advanced-pagination
  >
    <template #created="{ item }">
      {{ item.created | date }}
    </template>
    <template #rating="{ item }">
      <rating-field
        :value="item.rating"
        readonly
      />
    </template>
  </c-advanced-data-table>
</template>

<script>
import { localQueryMixin } from '@/mixins/query-local/query';
import { entitiesRemediationInstructionStatsMixin } from '@/mixins/entities/remediation/instruction-stats';

import RatingField from '@/components/forms/fields/rating-field.vue';

export default {
  components: { RatingField },
  mixins: [localQueryMixin, entitiesRemediationInstructionStatsMixin],
  props: {
    remediationInstruction: {
      type: Object,
      required: true,
    },
    interval: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      remediationInstructionStatsComments: [],
      totalItems: 0,
      pending: false,
    };
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.date'),
          value: 'created',
          sortable: false,
        },
        {
          text: this.$t('common.username'),
          value: 'user.display_name',
          sortable: false,
        },
        {
          text: this.$tc('common.rating'),
          value: 'rating',
          sortable: false,
        },
        {
          text: this.$tc('common.comment', 2),
          value: 'comment',
          sortable: false,
        },
      ];
    },
  },
  watch: {
    interval: 'fetchList',
  },
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const params = this.getQuery();

      params.from = this.interval.from;
      params.to = this.interval.to;

      const {
        data: remediationInstructionStatsComments,
        meta,
      } = await this.fetchRemediationInstructionStatsCommentsListWithoutStore({
        params,
        id: this.remediationInstruction._id,
      });
      this.remediationInstructionStatsComments = remediationInstructionStatsComments;
      this.totalItems = meta.total_count;

      this.pending = false;
    },
  },
};
</script>
