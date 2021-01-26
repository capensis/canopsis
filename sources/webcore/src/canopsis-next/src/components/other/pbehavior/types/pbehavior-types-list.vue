<template lang="pug">
  c-advanced-data-table.white(
    :headers="headers",
    :items="pbehaviorTypes",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledType",
    select-all,
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyPbehaviorTypeAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="icon_name", slot-scope="props")
      span.pbehavior-type-icon(v-if="props.item.icon_name")
        v-icon(color="white", size="18") {{ props.item.icon_name }}
    template(slot="actions", slot-scope="props")
      v-layout
        c-action-btn(
          :disabled="!props.item.editable",
          :tooltip="props.item.editable ? $t('common.edit') : $t('pbehaviorTypes.defaultType')",
          type="edit",
          @click="$emit('edit', props.item)"
        )
        c-action-btn(
          :disabled="!props.item.deletable",
          :tooltip="props.item.deletable ? $t('common.delete') : $t('pbehaviorTypes.defaultType')",
          type="delete",
          @click="$emit('remove', props.item._id)"
        )
    template(slot="expand", slot-scope="props")
      pbehavior-types-list-expand-panel(:pbehaviorType="props.item")
</template>

<script>
import rightsTechnicalPbehaviorTypesMixin from '@/mixins/rights/technical/pbehavior-types';

import PbehaviorTypesListExpandPanel from './partials/pbehavior-types-list-expand-panel.vue';

export default {
  components: {
    PbehaviorTypesListExpandPanel,
  },
  mixins: [rightsTechnicalPbehaviorTypesMixin],
  props: {
    pbehaviorTypes: {
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
          text: this.$t('common.icon'),
          value: 'icon_name',
          sortable: false,
        },
        {
          text: this.$t('common.priority'),
          value: 'priority',
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
    isDisabledType({ deletable }) {
      return !deletable;
    },
  },
};
</script>

<style lang="scss" scoped>
  .pbehavior-type-icon {
    display: inline-flex;
    padding: 2px 10px;
    border-radius: 10px;
    background: #17ffff;
  }
</style>
