<template lang="pug">
  advanced-data-table.white(
    :headers="headers",
    :items="pbehaviorExceptions",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    :is-disabled-item="isDisabledException",
    select-all,
    expand,
    search,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar", slot-scope="props")
      v-flex(v-show="hasDeleteAnyPbehaviorExceptionAccess && props.selected.length", xs4)
        v-btn(@click="$emit('remove-selected', props.selected)", icon)
          v-icon delete
    template(slot="actions", slot-scope="props")
      v-layout
        v-btn.mx-0(
          v-if="hasUpdateAnyPbehaviorExceptionAccess",
          icon,
          small,
          @click.stop="$emit('edit', props.item)"
        )
          v-icon edit
        v-tooltip(bottom, :disabled="props.item.deletable")
          v-btn.mx-0(
            slot="activator",
            v-if="hasDeleteAnyPbehaviorExceptionAccess",
            :disabled="!props.item.deletable",
            icon,
            small,
            @click.stop="$emit('remove', props.item._id)"
          )
            v-icon(color="error") delete
          span {{ $t('pbehaviorExceptions.usingException') }}
    template(slot="expand", slot-scope="props")
      pbehavior-exceptions-list-expand-panel(:pbehaviorException="props.item")
</template>

<script>
import rightsTechnicalPbehaviorExceptionsMixin from '@/mixins/rights/technical/pbehavior-exceptions';

import PbehaviorExceptionsListExpandPanel from './partials/pbehavior-exceptions-list-expand-panel.vue';

export default {
  components: {
    PbehaviorExceptionsListExpandPanel,
  },
  mixins: [rightsTechnicalPbehaviorExceptionsMixin],
  props: {
    pbehaviorExceptions: {
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
  },
  methods: {
    isDisabledException({ deletable }) {
      return !deletable;
    },
  },
};
</script>
