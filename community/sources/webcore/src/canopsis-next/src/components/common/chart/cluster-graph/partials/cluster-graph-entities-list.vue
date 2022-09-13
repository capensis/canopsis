<template lang="pug">
  div
    h6.my-2.title.text-xs-center {{ $tc('common.entity', 2) }}
    v-card.ma-2(v-for="(entity, index) in entities", :key="entity.key")
      v-card-text
        c-entity-field(
          :value="entity.data",
          :name="`entity-${entity.key}`",
          :entity-types="entityTypes",
          :item-disabled="isItemDisabled",
          :item-text="getItemText",
          return-object,
          clearable,
          @input="updateData($event, index)"
        )
        v-expand-transition
          v-combobox(
            v-if="entity.data",
            :value="entity.pinned",
            :items="pinnedListById[entity.data._id]",
            :loading="pinnedPendingById[entity.data._id]",
            :label="$t('modals.createTreeOfDependenciesMap.pinnedEntities')",
            :item-text="getItemText",
            item-value="_id",
            deletable-chips,
            chips,
            multiple,
            @change="updatePinned($event, index)"
          )
      v-card-actions
        v-layout(justify-end)
          c-action-btn(type="delete", @click="remove(index)")
    v-btn(color="primary", @click="add") {{ $t('modals.createTreeOfDependenciesMap.addEntity') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { ENTITY_TYPES, MAX_LIMIT } from '@/constants';

import uid from '@/helpers/uid';

import { formArrayMixin } from '@/mixins/form';

const { mapActions } = createNamespacedHelpers('service');

export default {
  mixins: [formArrayMixin],
  model: {
    prop: 'entities',
    event: 'input',
  },
  props: {
    entities: {
      type: Array,
      default: () => [],
    },
    impact: {
      type: Boolean,
      default: false,
    },
    entityTypes: {
      type: Array,
      default: () => [ENTITY_TYPES.service],
    },
  },
  data() {
    return {
      pinnedListById: {},
      pinnedPendingById: {},
    };
  },
  watch: {
    impact() {
      this.entities.forEach(({ data }) => {
        if (!data) {
          return;
        }

        this.$set(this.pinnedListById, data._id, []);

        this.fetchPinnedEntitiesList(data._id);
      });
    },
  },
  mounted() {
    if (this.entities.length) {
      this.entities.forEach(({ data }) => this.fetchPinnedEntitiesList(data._id));
    }
  },
  methods: {
    ...mapActions({
      fetchServiceDependenciesWithoutStore: 'fetchDependenciesWithoutStore',
      fetchServiceImpactsWithoutStore: 'fetchImpactsWithoutStore',
    }),

    getItemText(item) {
      return item.type === ENTITY_TYPES.service ? item.name : item._id;
    },

    isItemDisabled(item) {
      return this.entities.some(({ data }) => data?._id === item._id);
    },

    add() {
      const newEntity = { key: uid(), data: undefined, pinned: [] };

      this.addItemIntoArray(newEntity);

      this.$emit('add', newEntity);
    },

    remove(index) {
      const entity = this.entities[index];

      this.removeItemFromArray(index);

      this.$emit('remove', entity);
    },

    updateData(newData, index) {
      const oldEntity = this.entities[index] ?? {};
      const newEntity = { ...oldEntity, data: newData, pinned: [] };

      if (newData) {
        this.fetchPinnedEntitiesList(newData._id);
      }

      this.updateItemInArray(index, newEntity);

      this.$emit('update:data', newEntity, oldEntity);
    },

    updatePinned(newPinned, index) {
      const oldEntity = this.entities[index] ?? {};
      const newEntity = { ...oldEntity, pinned: newPinned };

      this.updateItemInArray(index, newEntity);

      this.$emit('update:pinned', newEntity, oldEntity);
    },

    async fetchPinnedEntitiesList(id, params = { limit: MAX_LIMIT }) {
      this.$set(this.pinnedPendingById, id, true);

      const { data } = await this.fetchDependenciesList({ id, params });

      this.$set(this.pinnedListById, id, data);
      this.$set(this.pinnedPendingById, id, false);
    },

    fetchDependenciesList(data) {
      return this.impact
        ? this.fetchServiceImpactsWithoutStore(data)
        : this.fetchServiceDependenciesWithoutStore(data);
    },
  },
};
</script>
