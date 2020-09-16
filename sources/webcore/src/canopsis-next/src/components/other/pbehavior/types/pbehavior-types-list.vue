<template lang="pug">
  advanced-data-table.white(
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
        v-btn(@click="deleteTypes(props.selected)", icon)
          v-icon delete
    template(slot="icon_name", slot-scope="props")
      span.pbehavior-type-icon(v-if="props.item.icon_name")
        v-icon(color="white", size="18") {{ props.item.icon_name }}
    template(slot="actions", slot-scope="props")
      v-layout
        v-tooltip(bottom, :disabled="props.item.editable")
          v-btn.mx-0(
            slot="activator",
            v-if="hasUpdateAnyPbehaviorTypeAccess",
            :disabled="!props.item.editable",
            icon,
            small,
            @click.stop="$emit('edit', props.item)"
          )
            v-icon edit
          span {{ $t('pbehaviorTypes.defaultType') }}
        v-tooltip(bottom, :disabled="props.item.deletable")
          v-btn.mx-0(
            slot="activator",
            v-if="hasDeleteAnyPbehaviorTypeAccess",
            :disabled="!props.item.deletable",
            icon,
            small,
            @click.stop="$emit('remove', props.item._id)"
          )
            v-icon(color="error") delete
          span {{ $t('pbehaviorTypes.usingType') }}
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

    deleteTypes(selected) {
      this.$emit('remove-selected', selected);
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
