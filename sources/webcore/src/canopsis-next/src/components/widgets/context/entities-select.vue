<template lang="pug">
  choose-expansion-panel(
    :entities="entities",
    :label="label",
    @remove="removeEntity"
  )
    v-btn.error(v-show="entities.length", small, @click="clear") {{ $t('common.clear') }}
    context-general-list(@update:selectedIds="updateEntities($event)")
</template>

<script>
import { union, filter } from 'lodash';

import formBaseMixin from '@/mixins/form/base';

import ContextGeneralList from '@/components/widgets/context/context-general-list.vue';
import ChooseExpansionPanel from '@/components/common/choose-expansion-panel/choose-expansion-panel.vue';

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
  components: { ChooseExpansionPanel, ContextGeneralList },
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
      default: () => [],
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
