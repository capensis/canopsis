<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="remediationJobs",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledJob",
    :select-all="removable",
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-if="removable",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#actions="{ item, disabled }")
      c-action-btn(
        v-if="updatable",
        type="edit",
        @click="$emit('edit', item)"
      )
      c-action-btn(
        v-if="duplicable",
        type="duplicate",
        @click="$emit('duplicate', item)"
      )
      c-action-btn(
        v-if="removable",
        :tooltip="disabled ? $t('remediation.job.usingJob') : $t('common.delete')",
        :disabled="disabled",
        type="delete",
        @click="$emit('remove', item)"
      )
</template>

<script>
export default {
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
    removable: {
      type: Boolean,
      default: false,
    },
    updatable: {
      type: Boolean,
      default: false,
    },
    duplicable: {
      type: Boolean,
      default: false,
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
          value: 'author.name',
        },
        {
          text: this.$t('remediation.job.configuration'),
          value: 'config.name',
          sortable: false,
        },
        {
          text: this.$t('remediation.job.jobId'),
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
