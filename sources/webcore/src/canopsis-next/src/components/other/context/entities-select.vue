<template lang="pug">
  v-expansion-panel.my-1
    v-expansion-panel-content.grey.darken-2.white--text
      div.white--text(slot="header") {{ label }}
      v-card
        v-layout(wrap)
          v-chip(
            v-for="entity in entities",
            :key="entity._id",
            close,
            @input="removeEntity(entity)"
          ) {{ entity }}
        v-btn.red.white--text(v-show="entities.length", @click="clear", small) Clear
        context-general-list(@update:selectedIds="updateEntities($event)")
</template>

<script>
import { union, filter } from 'lodash';

import formBaseMixin from '@/mixins/form/base';

import ContextGeneralList from '@/components/other/context/context-general-list.vue';

/**
 * Component to select entities for impact/dependencies
 *
 * @module context
 *
 * @prop {String} [label] - "Impacts" or "Dependencies"
 *
 * @event selectedIds#update
 */
export default {
  components: { ContextGeneralList },
  mixins: [formBaseMixin],
  model: {
    prop: 'entities',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    entities: {
      type: Array,
      default() {
        return [];
      },
    },
  },
  methods: {
    updateEntities(entities) {
      const entitiesIds = entities.map(entity => entity._id);
      const selectedEntities = union(entitiesIds, this.entities);

      this.updateModel(selectedEntities);
    },
    clear() {
      this.updateModel([]);
    },
    removeEntity(entity) {
      const updatedEntities = filter(this.entities, item => item !== entity);

      this.updateModel(updatedEntities);
    },
  },
};
</script>

<style scoped>
  .addContainer {
    display: flex;
    align-items: center;
    align-content: flex-start;
  }

  .entities {
    width: 50%;
    white-space: nowrap;
    overflow: auto;
  }

  .label {
    width: 16%;
  }

  .scrollbar::-webkit-scrollbar-track {
    box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
    -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.3);
    border-radius: 10px;
    background-color: #F5F5F5;
  }

  .scrollbar::-webkit-scrollbar {
    height: 3px;
  }

  .scrollbar::-webkit-scrollbar-thumb {
    border-radius: 10px;
    box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
    -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, .3);
    background-color: darkgray;
    height: 0;
  }
</style>
