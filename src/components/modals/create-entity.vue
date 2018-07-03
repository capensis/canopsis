<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createEntity.title') }}
    v-container
      v-card-text
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
        v-layout(row, justify-space-around, align-center)
          v-flex(xs6)
            v-select(
              :items="types",
              v-model="form.type",
              label="Type",
              single-line
            )
          v-spacer
          v-flex
            v-switch(:label="$t('common.enabled')", v-model="enabled", hide-details)
      entities-select.my-1(label="Impacts", :entities.sync=entities)
      entities-select.my-1(label="Dependencies", :entities.sync=entities)
    v-card-actions.mt-2
      v-btn(@click.prevent="submit", color="green darken-4 white--text") {{ $t('common.submit') }}
      v-btn(@click.prevent="manageInfos", color="green darken-4 white--text") {{ $t('modals.createEntity.fields.manageInfos') }}
</template>

<script>
import { MODALS } from '@/constants';
import EntitiesSelect from '@/components/other/context/actions/create-entities/entities-select.vue';

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
