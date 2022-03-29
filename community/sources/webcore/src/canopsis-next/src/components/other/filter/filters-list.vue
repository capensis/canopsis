<template lang="pug">
  div
    v-alert(
      :value="!filters.length",
      type="info"
    ) {{ $t('modals.createFilter.emptyFilters') }}
    filter-tile(
      v-for="filter in filters",
      :filter="filter",
      :key="filter._id",
      :editable="editable",
      @edit="showEditFilterModal(filter)",
      @delete="showDeleteFilterModal(filter)"
    )
    v-btn.ml-0(
      v-if="addable",
      color="primary",
      outline,
      @click.prevent="showCreateFilterModal"
    ) {{ $t('common.addFilter') }}
</template>

<script>
import { MODALS } from '@/constants';

import { entitiesWidgetMixin } from '@/mixins/entities/view/widget';

import FilterTile from './partials/filter-tile.vue';

export default {
  components: { FilterTile },
  mixins: [entitiesWidgetMixin],
  props: {
    widgetId: {
      type: String,
      required: true,
    },
    private: {
      type: Boolean,
      default: false,
    },
    addable: {
      type: Boolean,
      default: true,
    },
    editable: {
      type: Boolean,
      default: true,
    },
    withAlarm: {
      type: Boolean,
      default: false,
    },
    withEntity: {
      type: Boolean,
      default: false,
    },
    withPbehavior: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    widget() {
      return this.getWidgetById(this.widgetId);
    },

    filters() {
      return (this.widget?.filters ?? []).filter(filter => filter.is_private === this.private);
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.createFilter.create.title'),
          withTitle: true,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          action: () => {},
        },
      });
    },

    showEditFilterModal(filter) {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          filter,

          title: this.$t('modals.createFilter.edit.title'),
          withTitle: true,
          withAlarm: this.withAlarm,
          withEntity: this.withEntity,
          withPbehavior: this.withPbehavior,
          action: () => {},
        },
      });
    },

    showDeleteFilterModal() {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => {},
        },
      });
    },
  },
};
</script>
