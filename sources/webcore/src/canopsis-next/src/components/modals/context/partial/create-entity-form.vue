<template lang="pug">
  div(data-test="createEntityForm")
    v-container(fluid)
      v-layout(row)
        v-text-field(
          data-test="entityFormName",
          v-field="form.name",
          v-validate="'required'",
          :label="$t('common.name')",
          :value="form.name",
          :error-messages="errors.collect('name')",
          name="name"
        )
      v-layout(row)
        v-textarea(
          data-test="entityFormDescription",
          v-field="form.description",
          :label="$t('common.description')"
        )
      v-layout(data-test="entityFormFieldLayout", row)
        v-switch(
          data-test="entityFormEnabled",
          v-field="form.enabled",
          :label="$t('common.enabled')",
          color="primary"
        )
        v-select(
          v-field="form.type",
          v-validate="'required'",
          :items="types",
          :error-messages="errors.collect('type')",
          :label="$t('modals.createEntity.fields.type')",
          name="type",
          single-line
        )
      v-layout(wrap)
        v-flex(data-test="entityFormImpact", xs12)
          entities-select(
            v-field="form.impact",
            :label="$t('modals.createEntity.fields.impact')"
          )
        v-flex(data-test="entityFormDepends", xs12)
          entities-select(
            v-field="form.depends",
            :label="$t('modals.createEntity.fields.depends')"
          )
</template>

<script>
import EntitiesSelect from '../entities-select.vue';

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
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
