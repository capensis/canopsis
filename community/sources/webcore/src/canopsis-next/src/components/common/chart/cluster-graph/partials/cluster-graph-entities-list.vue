<template>
  <div>
    <h6 class="my-2 text-h6 text-center">
      {{ $tc('common.entity', 2) }}
    </h6>
    <v-card
      v-for="(entityItem, index) in entities"
      :key="entityItem.key"
      class="ma-2"
    >
      <v-card-text>
        <c-entity-field
          :value="entityItem.entity"
          :name="`entity-${entityItem.key}`"
          :entity-types="entityTypes"
          :item-disabled="isItemDisabled"
          :item-text="getItemText"
          required
          return-object
          clearable
          autocomplete
          @input="updateEntity($event, index)"
        />
        <v-expand-transition>
          <v-combobox
            v-if="entityItem.entity"
            :value="entityItem.pinned"
            :items="pinnedListById[entityItem.entity._id]"
            :loading="pinnedPendingById[entityItem.entity._id]"
            :label="$t('modals.createTreeOfDependenciesMap.pinnedEntities')"
            :item-text="getItemText"
            item-value="_id"
            deletable-chips
            chips
            multiple
            @change="updatePinned($event, index)"
          />
        </v-expand-transition>
      </v-card-text>
      <v-card-actions>
        <v-layout justify-end>
          <c-action-btn
            :disabled="entities.length === 1"
            type="delete"
            @click="remove(index)"
          />
        </v-layout>
      </v-card-actions>
    </v-card>
    <v-btn
      class="ml-2"
      color="primary"
      @click="add"
    >
      {{ $t('modals.createTreeOfDependenciesMap.addEntity') }}
    </v-btn>
  </div>
</template>

<script>
import { ENTITY_TYPES, MAX_LIMIT } from '@/constants';

import { uid } from '@/helpers/uid';
import { getMapEntityText } from '@/helpers/entities/map/list';

import { formArrayMixin } from '@/mixins/form';
import { entitiesEntityDependenciesMixin } from '@/mixins/entities/entity-dependencies';

export default {
  mixins: [formArrayMixin, entitiesEntityDependenciesMixin],
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
  computed: {
    notEmptyEntities() {
      return this.entities.filter(({ entity }) => entity);
    },
  },
  watch: {
    impact() {
      this.entities.forEach(({ entity }) => {
        if (!entity) {
          return;
        }

        this.$set(this.pinnedListById, entity._id, []);

        this.fetchPinnedEntitiesList(entity._id);
      });
    },
  },
  mounted() {
    if (this.notEmptyEntities.length) {
      this.notEmptyEntities.forEach(({ entity }) => this.fetchPinnedEntitiesList(entity._id));
    }
  },
  methods: {
    getItemText(item) {
      return getMapEntityText(item);
    },

    isItemDisabled(item) {
      return this.entities.some(({ entity }) => entity?._id === item._id);
    },

    add() {
      const newEntity = { key: uid(), entity: undefined, pinned: [] };

      this.addItemIntoArray(newEntity);

      this.$emit('add', newEntity);
    },

    remove(index) {
      const entity = this.entities[index];

      this.removeItemFromArray(index);

      this.$emit('remove', entity);
    },

    updateEntity(newEntity, index) {
      const oldEntityItem = this.entities[index] ?? {};
      const newEntityItem = { ...oldEntityItem, entity: newEntity, pinned: [] };

      if (newEntity) {
        this.fetchPinnedEntitiesList(newEntity._id);
      }

      this.updateItemInArray(index, newEntityItem);

      this.$emit('update:entity', newEntityItem, oldEntityItem);
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
  },
};
</script>
