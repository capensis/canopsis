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
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyPbehaviorReasonAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="actions", slot-scope="props")
      c-action-btn(
        v-if="hasUpdateAnyPbehaviorReasonAccess",
        type="edit",
        @click="$emit('edit', props.item)"
      )
      c-action-btn(
        v-if="hasDeleteAnyPbehaviorReasonAccess",
        :tooltip="props.item.deletable ? $t('common.delete') : $t('pbehaviorReasons.usingReason')",
        :disabled="!props.item.deletable",
        type="delete",
        @click="$emit('remove', props.item._id)"
      )
    template(slot="expand", slot-scope="props")
      pbehavior-reasons-list-expand-panel(:pbehaviorReason="props.item")
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
