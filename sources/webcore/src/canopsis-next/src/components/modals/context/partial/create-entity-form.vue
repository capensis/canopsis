<template lang="pug">
  div
    v-container(fluid)
      v-layout(row)
        v-text-field(
        :label="$t('common.name')",
        :value="form.name",
        @input="updateField('name', $event)",
        :error-messages="errors.collect('name')",
        v-validate="'required'",
        data-vv-name="name"
        )
      v-layout(row)
        v-textarea(
        :label="$t('common.description')",
        :value="form.description",
        @input="updateField('description', $event)",
        data-vv-name="description",
        :error-messages="errors.collect('description')"
        )
      v-layout(row)
        v-switch(
        color="primary",
        :label="$t('common.enabled')",
        :input-value="form.enabled",
        @change="updateField('enabled', $event)"
        )
        v-select(
        :items="types",
        :value="entityType",
        data-vv-name="type",
        v-validate="'required'",
        :error-messages="errors.collect('type')",
        @input="updateField('type', $event)",
        :label="$t('modals.createEntity.fields.type')",
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
  data() {
    return {
      showValidationErrors: true,
      types: [
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
      ],
    };
  },
  computed: {
    entityType() {
      let entityType;
      this.types.map((item, index) => {
        if (this.form.type === item.value) {
          return entityType = this.types[index].value;
        }
        return null;
      });
      return entityType;
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
