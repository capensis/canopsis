<template lang="pug">
  c-advanced-data-table(
    :items="remediationInstructionStatsComments",
    :headers="headers",
    :loading="pending",
    :pagination.sync="pagination",
    :total-items="totalItems",
    advanced-pagination
  )
    template(slot="created", slot-scope="props") {{ props.item.created | date('long', true) }}
    template(slot="rating", slot-scope="props")
      rating-field(:value="props.item.rating", readonly)
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
          value: 'user.name',
          sortable: false,
        },
        {
          text: this.$t('remediationInstructionStats.rating'),
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
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const {
        data: remediationInstructionStatsComments,
        meta,
      } = await this.fetchRemediationInstructionStatsCommentsListWithoutStore({
        id: this.remediationInstruction._id,
        params: this.getQuery(),
      });
      this.remediationInstructionStatsComments = remediationInstructionStatsComments;
      this.totalItems = meta.total_count;

      this.pending = false;
    },
  },
};
</script>