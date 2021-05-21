<template lang="pug">
  c-advanced-data-table.white(
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
      c-action-btn(
        v-if="hasUpdateAnyRemediationJobAccess",
        type="edit",
        @click="$emit('edit', props.item)"
      )
      c-action-btn(
        v-if="hasDeleteAnyRemediationJobAccess",
        :tooltip="props.disabled ? $t('remediationJobs.usingJob') : $t('common.delete')",
        :disabled="props.disabled",
        type="delete",
        @click="$emit('remove', props.item)"
      )
</template>

<script>
import { permissionsTechnicalRemediationJobMixin } from '@/mixins/permissions/technical/remediation-job';

export default {
  mixins: [permissionsTechnicalRemediationJobMixin],
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
          value: 'author.username',
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
