<template lang="pug">
  div
    v-container(fluid)
      v-layout(row)
        v-text-field(
          v-field="form.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :error-messages="errors.collect('name')",
          name="name"
        )
      v-layout(row)
        v-textarea(
          v-field="form.description",
          :label="$t('common.description')",
          :error-messages="errors.collect('description')",
          name="description"
        )
      v-layout(row)
        v-switch(
          v-field="form.enabled",
          :label="$t('common.enabled')",
          color="primary"
        )
        v-select(
          v-model="form.type",
          v-validate="'required'",
          :items="types",
          :error-messages="errors.collect('type')",
          :label="$t('modals.createEntity.fields.type')",
          name="type",
          single-line
        )
      v-layout(wrap)
        v-flex(xs12)
          entities-select(
            :label="$t('modals.createEntity.fields.impact')",
            :entities="form.impact",
            @updateEntities="updateImpact"
          )
        v-flex(xs12)
          entities-select(
            :label="$t('modals.createEntity.fields.depends')",
            :entities="form.depends",
            @updateEntities="updateDepends"
          )
</template>

<script>
import formMixin from '@/mixins/form';

import EntitiesSelect from '@/components/other/context/entities-select.vue';

/**
 * Form to create a new entity
 *
 * @prop {String} [name] - Name of the entity (null if creating a new entity)
 * @prop {String} [description] - Description on the entity (null if creating a new entity)
 * @prop {String} [type] - Type of the entity (null if creating a new entity)
 * @prop {Array} [impact] - List of the entity's impacts (null if creating a new entity)
 * @prop {Array} [depends] - List of the entity's depends (null if creating a new entity)
 * @prop {Boolean} [enabled] - Whether the entity is enabled or not
 *
 * @event name#update
 * @event description#update
 * @event enabled#update
 * @event type#update
 *
 * @module context
 */
export default {
  inject: ['$validator'],
  components: {
    EntitiesSelect,
  },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  computed: {
    types() {
      return [
        {
          text: this.$t('modals.createEntity.fields.types.connector'),
          value: 'connector',
        },
        {
          text: this.$t('modals.createEntity.fields.types.component'),
          value: 'component',
        },
        {
          text: this.$t('modals.createEntity.fields.types.resource'),
          value: 'resource',
        },
      ];
    },
  },
  methods: {
    updateImpact(entities) {
      this.updateField('impact', entities);
    },
    updateDepends(entities) {
      this.updateField('depends', entities);
    },
  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
