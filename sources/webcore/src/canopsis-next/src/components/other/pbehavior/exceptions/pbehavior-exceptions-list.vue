<template lang="pug">
  advanced-data-table.white(
    v-model="selected",
    :headers="headers",
    :items="pbehaviorExceptions",
    :loading="pending",
    :total-items="totalItems",
    :pagination="pagination",
    select-all,
    expand,
    search,
    hide-actions,
    advanced-pagination,
    @update:pagination="$emit('update:pagination', $event)"
  )
    template(slot="toolbar")
      v-flex(v-show="hasDeleteAnyPbehaviorExceptionAccess && selectedExceptions.length", xs4)
        v-btn(@click="deleteSelectedExceptions", icon)
          v-icon delete
    template(slot="selectAll", slot-scope="props")
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

    selectedExceptions() {
      return this.selected.filter(({ deletable }) => deletable);
    },
  },
  methods: {
    deleteSelectedExceptions() {
      this.$emit('remove-selected', this.selectedExceptions);
    },
  },
};
</script>
