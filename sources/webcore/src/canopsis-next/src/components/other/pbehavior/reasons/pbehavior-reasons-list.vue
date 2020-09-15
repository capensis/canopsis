<template lang="pug">
  advanced-data-table.white(
    v-model="selected",
    :headers="headers",
    :items="pbehaviorReasons",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    select-all,
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar")
      v-flex(v-show="hasDeleteAnyPbehaviorReasonAccess && selectedReasons.length", xs4)
        v-btn(@click="deleteSelectedReasons", icon)
          v-icon delete
    template(slot="item-select", slot-scope="props")
      v-checkbox-functional(
        v-if="props.item.deletable",
        :inputValue="props.selected",
        primary,
        hide-details,
        @change="props.select"
      )
      v-checkbox-functional(v-else, disabled, primary, hide-details)
    template(slot="actions", slot-scope="props")
      v-layout
        v-btn.mx-0(
          slot="activator",
          v-if="hasUpdateAnyPbehaviorReasonAccess",
          icon,
          small,
          @click.stop="$emit('edit', props.item)"
        )
          v-icon edit
        v-tooltip(bottom, :disabled="props.item.deletable")
          v-btn.mx-0(
            slot="activator",
            v-if="hasDeleteAnyPbehaviorReasonAccess",
            :disabled="!props.item.deletable",
            icon,
            small,
            @click.stop="$emit('remove', props.item._id)"
          )
            v-icon(color="error") delete
          span {{ $t('pbehaviorReasons.usingReason') }}
    template(slot="expand", slot-scope="props")
      pbehavior-reasons-list-expand-panel(:pbehaviorReason="props.item")
</template>

<script>
import rightsTechnicalPbehaviorReasonsMixin from '@/mixins/rights/technical/pbehavior-reasons';

import PbehaviorReasonsListExpandPanel from './partials/pbehavior-reasons-list-expand-panel.vue';

export default {
  components: {
    PbehaviorReasonsListExpandPanel,
  },
  mixins: [rightsTechnicalPbehaviorReasonsMixin],
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
  data() {
    return {
      selected: [],
    };
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

    selectedReasons() {
      return this.selected.filter(({ deletable }) => deletable);
    },
  },
  methods: {
    deleteSelectedReasons() {
      this.$emit('remove-selected', this.selectedReasons);
    },
  },
};
</script>
