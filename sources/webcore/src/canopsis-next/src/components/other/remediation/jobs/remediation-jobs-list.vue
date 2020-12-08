<template lang="pug">
  advanced-data-table.white(
    :headers="headers",
    :items="remediationJobs",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledJob",
    select-all,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyRemediationJobAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="actions", slot-scope="props")
      v-layout
        v-btn.mx-0(
          v-if="hasUpdateAnyRemediationJobAccess",
          icon,
          small,
          @click.stop="$emit('edit', props.item)"
        )
          v-icon edit
        v-tooltip(bottom, :disabled="!props.disabled")
          v-btn.mx-0(
            slot="activator",
            v-if="hasDeleteAnyRemediationJobAccess",
            :disabled="props.disabled",
            icon,
            small,
            @click.stop="$emit('remove', props.item)"
          )
            v-icon(color="error") delete
          span {{ $t('remediationJobs.usingJob') }}
</template>

<script>
import rightsTechnicalRemediationJobMixin from '@/mixins/rights/technical/remediation-job';

export default {
  mixins: [rightsTechnicalRemediationJobMixin],
  props: {
    remediationJobs: {
      type: Array,
      required: true,
    },
    pending: {
      type: Boolean,
      default: false,
    },
    totalItems: {
      type: Number,
      required: false,
    },
    pagination: {
      type: Object,
      required: true,
    },
  },
  computed: {
    headers() {
      return [
        {
          text: this.$t('common.name'),
          value: 'name',
        },
        {
          text: this.$t('common.author'),
          value: 'author',
        },
        {
          text: this.$t('remediationJobs.table.configuration'),
          value: 'config.name',
        },
        {
          text: this.$t('remediationJobs.table.jobId'),
          value: 'job_id',
        },
        {
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    isDisabledJob({ deletable }) {
      return !deletable;
    },
  },
};
</script>
