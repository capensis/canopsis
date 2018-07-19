<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createEntity.title') }}
    v-tabs
      v-tab(
        v-for="tab in tabs",
        :key="tab.name",
        @click.prevent="currentComponent = tab.component",
      ) {{ tab.name }}
      v-tab-item
        keep-alive
        create-form(
          :name.sync="form.name",
          :description.sync="form.description",
          :type.sync="form.type",
          :enabled.sync="form.enabled",
          :infos.sync="form.infos",
        )
      v-tab-item
        manage-infos(:infos.sync="form.infos")
    v-card-actions
      v-btn(@click.prevent="submit", color="blue darken-4 white--text") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import modalInnerMixin from '@/mixins/modal/modal-inner';
import { MODALS } from '@/constants';

import CreateForm from './create-entity-form.vue';

const { mapActions: entitiesMapActions } = createNamespacedHelpers('entity');

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: { CreateForm },
  mixins: [modalInnerMixin],
  data() {
    return {
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
      showValidationErrors: true,
      enabled: true,
      form: {
        name: '',
        description: '',
        type: '',
        enabled: true,
        depends: [],
        impact: [],
        infos: {},
      },
    };
  },
  mounted() {
    if (this.config.item) {
      this.form = { ...this.config.item };
    }
  },
  methods: {
    ...entitiesMapActions({
      createEntity: 'create',
      editEntity: 'edit',
    }),
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();
      if (formIsValid) {
        // If there's an item, means we're editing. If there's not, we're creating an entity
        if (this.config.item) {
          this.editEntity({ data: this.form });
        } else {
          const formData = { ...this.form, _id: this.form.name };
          this.createEntity({ data: formData });
        }
        this.hideModal();
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
