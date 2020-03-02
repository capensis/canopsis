<template lang="pug">
  div
    draggable(
      :value="filters",
      :options="draggableOptions",
      element="v-list",
      @change="changeFiltersOrdering"
    )
      v-list-tile.filter-item.pa-0(v-for="(filter, index) in filters", :key="filter.title")
        v-layout(:data-test="`filterItem-${filter.title}`")
          v-flex(xs12)
            v-list-tile-content {{ filter.title }}
          v-list-tile-action(v-if="hasAccessToEditFilter")
            v-layout
              v-btn.ma-1(
                icon,
                :data-test="`editFilter-${filter.title}`",
                @click="showEditFilterModal(index)"
              )
                v-icon edit
              v-btn.ma-1(
                icon,
                :data-test="`deleteFilter-${filter.title}`",
                @click="showDeleteFilterModal(index)"
              )
                v-icon delete
    v-btn.ml-0(
      data-test="addFilter",
      v-if="hasAccessToAddFilter",
      color="primary",
      @click.prevent="showCreateFilterModal"
    ) {{ $t('common.add') }}
</template>

<script>
import Draggable from 'vuedraggable';

import { VUETIFY_ANIMATION_DELAY } from '@/config';
import { MODALS, ENTITIES_TYPES } from '@/constants';

export default {
  components: { Draggable },
  props: {
    filters: {
      type: Array,
      default: () => [],
    },
    hasAccessToAddFilter: {
      type: Boolean,
      default: true,
    },
    hasAccessToEditFilter: {
      type: Boolean,
      default: true,
    },
    entitiesType: {
      type: String,
      default: ENTITIES_TYPES.alarm,
      validator: value => [ENTITIES_TYPES.alarm, ENTITIES_TYPES.entity].includes(value),
    },
  },
  computed: {
    existingTitles() {
      return this.filters.map(({ title }) => title);
    },

    draggableOptions() {
      return {
        animation: VUETIFY_ANIMATION_DELAY,
      };
    },
  },
  methods: {
    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.create.title'),
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
          action: newFilter => this.$emit('create:filter', newFilter),
        },
      });
    },

    showEditFilterModal(index) {
      const filter = this.filters[index];

      this.$modals.show({
        name: MODALS.createFilter,
        config: {
          title: this.$t('modals.filter.edit.title'),
          filter,
          entitiesType: this.entitiesType,
          existingTitles: this.existingTitles,
          action: newFilter => this.$emit('update:filter', newFilter, index),
        },
      });
    },

    showDeleteFilterModal(index) {
      this.$modals.show({
        name: MODALS.confirmation,
        config: {
          action: () => this.$emit('delete:filter', index),
        },
      });
    },

    changeFiltersOrdering({ moved, added, removed }) {
      const filters = [...this.filters];

      if (moved) {
        const [item] = filters.splice(moved.oldIndex, 1);

        filters.splice(moved.newIndex, 0, item);
      } else if (added) {
        filters.splice(added.newIndex, 0, added.element);
      } else if (removed) {
        filters.splice(removed.oldIndex, 1);
      }

      if (filters) {
        this.$emit('update:filters', filters);
      }
    },
  },
};
</script>

<style lang="scss" scoped>
  .filter-item {
    cursor: move;
  }
</style>
