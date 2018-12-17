<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t(config.title) }}
    v-tabs(slider-color="primary")
      v-tab(
      v-for="tab in tabs",
      :key="tab.name",
      @click.prevent="currentComponent = tab.component",
      ) {{ tab.name }}
      v-tab-item
        keep-alive
        create-form(v-model="form")
      v-tab-item
        manage-infos(v-model="form.infos")
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat, v-if="!submitting") {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import uuid from '@/helpers/uuid';

import popupMixin from '@/mixins/popup';
import modalInnerMixin from '@/mixins/modal/inner';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

import CreateForm from './partial/create-entity-form.vue';
import ManageInfos from './partial/manage-infos.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CreateForm,
    ManageInfos,
  },
  mixins: [
    popupMixin,
    modalInnerMixin,
    entitiesContextEntityMixin,
  ],
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
      tabs: [
        { component: 'CreateForm', name: this.$t('modals.createEntity.fields.form') },
        { component: 'ManageInfos', name: this.$t('modals.createEntity.fields.manageInfos') },
      ],
      showValidationErrors: true,
      form: {
        name: '',
        description: '',
        type: '',
        enabled: true,
        depends: [],
        impact: [],
        infos: {},
      },
      submitting: false,
    };
  },
  mounted() {
    if (this.config.item) {
      this.form = { ...this.config.item };
    }

    if (this.config.isDuplicating) {
      this.form.name = '';
    }
  },
  methods: {
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async submit() {
      this.submitting = true;
      const formIsValid = await this.$validator.validateAll();
      if (formIsValid) {
        const formData = { ...this.form };

        if (!this.config.item || this.config.isDuplicating) {
          formData._id = uuid('entity');
        }
        try {
          await this.config.action(formData);

          this.refreshContextEntitiesLists();

          this.hideModal();
        } catch (err) {
          console.error(err);
          this.addErrorPopup({ text: this.$t('error.default') });
        }
      }

      this.submitting = false;
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
