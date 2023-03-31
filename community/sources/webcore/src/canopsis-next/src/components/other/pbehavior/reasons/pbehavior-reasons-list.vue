<template lang="pug">
  c-advanced-data-table(
    :headers="headers",
    :items="pbehaviorReasons",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledReason",
    select-all,
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(#mass-actions="{ selected }")
      c-action-btn(
        v-if="hasDeleteAnyPbehaviorReasonAccess",
        type="delete",
        @click="$emit('remove-selected', selected)"
      )
    template(#actions="{ item }")
      c-action-btn(
        v-if="hasUpdateAnyPbehaviorReasonAccess",
        type="edit",
        @click="$emit('edit', item)"
      )
      c-action-btn(
        v-if="hasDeleteAnyPbehaviorReasonAccess",
        :tooltip="item.deletable ? $t('common.delete') : $t('pbehavior.reasons.usingReason')",
        :disabled="!item.deletable",
        type="delete",
        @click="$emit('remove', item._id)"
      )
    template(#expand="{ item }")
      pbehavior-reasons-list-expand-panel(:pbehavior-reason="item")
</template>

<script>
import { permissionsTechnicalPbehaviorReasonsMixin } from '@/mixins/permissions/technical/pbehavior-reasons';

import PbehaviorReasonsListExpandPanel from './partials/pbehavior-reasons-list-expand-panel.vue';

export default {
  components: {
    PbehaviorReasonsListExpandPanel,
  },
  mixins: [permissionsTechnicalPbehaviorReasonsMixin],
  props: {
    pbehaviorReasons: {
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
          text: this.$t('common.actionsLabel'),
          value: 'actions',
          sortable: false,
        },
      ];
    },
  },
  methods: {
    isDisabledReason({ deletable }) {
      return !deletable;
    },
  },
};
</script>
