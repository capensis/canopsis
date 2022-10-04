<template lang="pug">
  choose-expansion-panel(
    :entities="entitiesIds",
    :disabled-entities="disabledEntitiesIds",
    :existing-entities="existingEntitiesIds",
    :label="label",
    :errors="errors.collect(name)",
    clearable,
    @clear="clear",
    @remove="removeEntity"
  )
    choose-entities-list(:entities-ids="entitiesIds", @select="updateEntities")
</template>

<script>
import { formBaseMixin } from '@/mixins/form';

import ChooseEntitiesList from '@/components/widgets/context/choose-entities-list.vue';
import ChooseExpansionPanel from '@/components/common/choose-expansion-panel/choose-expansion-panel.vue';

export default {
  inject: ['$validator'],
  components: { ChooseExpansionPanel, ChooseEntitiesList },
  mixins: [formBaseMixin],
  model: {
    prop: 'entitiesIds',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      required: true,
    },
    entitiesIds: {
      type: Array,
      default: () => [],
    },
    disabledEntitiesIds: {
      type: Array,
      default: () => [],
    },
    existingEntitiesIds: {
      type: Array,
      default: () => [],
    },
    name: {
      type: String,
      default: 'entities',
    },
  },
  watch: {
    existingEntitiesIds: 'validate',
    entitiesIds: 'validate',
  },
  created() {
    this.$validator.attach({
      name: this.name,
      rules: {
        unique: {
          values: this.existingEntitiesIds,
        },
      },
      getter: () => this.entitiesIds,
      vm: this,
    });
  },
  beforeDestroy() {
    this.$validator.detach(this.name);
  },
  methods: {
    validate() {
      this.$validator.validate(this.name);
    },

    clear() {
      this.updateModel([]);
    },

    updateEntities(entities) {
      const entityIds = entities.map(({ _id }) => _id);

      this.updateModel([...this.entitiesIds, ...entityIds]);
    },

    removeEntity(entityId) {
      const updatedEntities = this.entitiesIds.filter(id => id !== entityId);

      this.updateModel(updatedEntities);
    },
  },
};
</script>
