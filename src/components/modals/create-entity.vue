<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createEntity.title') }}
    v-card-text
      v-container
        v-layout(row)
          v-text-field(
          :label="$t('common.name')",
          :error-messages="errors.collect('ticket')",
          v-model="form.name",
          v-validate="'required'",
          data-vv-name="name"
          )
        v-layout(row)
          v-text-field(
          :label="$t('common.description')",
          :error-messages="errors.collect('output')",
          v-model="form.description",
          v-validate="'required'",
          data-vv-name="description",
          multi-line
          )
        v-layout(row)
          v-switch(:label="$t('common.enabled')", v-model="enabled")
          v-select(
          :items="types"
          v-model="form.type"
          label="Type"
          single-line
          )
    entities-select(label="Impacts", @update:entities="updateImpact($event)")
    entities-select(label="Dependencies", @update:entities="updateDependencies($event)")
    v-card-actions
      v-btn(@click.prevent="submit", color="primary") {{ $t('common.submit') }}
      v-btn(@click.prevent="manageInfos", color="primary") {{ $t('modals.createEntity.fields.manageInfos') }}
</template>

<script>

import { MODALS } from '@/constants';
import EntitiesSelect from '@/components/other/context-explorer/actions/create-entities/entities-select.vue';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EntitiesSelect },
  data() {
    return {
      showValidationErrors: true,
      enabled: true,
      types: [
        this.$t('modals.createEntity.fields.types.connector'),
        this.$t('modals.createEntity.fields.types.component'),
        this.$t('modals.createEntity.fields.types.resource'),
      ],
      form: {
        name: '',
        description: '',
        type: '',
        impacts: [],
        dependencies: [],
      },
    };
  },
  methods: {
    async create() {
      // TO DO
      // Entity creation
    },
    async manageInfos() {
      // TO DO
      // manage infos
    },
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();
      if (formIsValid) {
        await this.create(true);
      }
    },

  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
  .impact {
    background-color: grey;
  }
</style>
