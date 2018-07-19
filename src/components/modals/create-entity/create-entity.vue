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
      v-btn(@click.prevent="submit", color="primary") {{ $t('common.submit') }}
</template>

<script>
import modalMixin from '@/mixins/modal/modal';
import modalInnerMixin from '@/mixins/modal/modal-inner';

import { MODALS } from '@/constants';

import CreateForm from './create-entity-form.vue';
import ManageInfos from './manage-infos.vue';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CreateForm,
    ManageInfos,
  },
  mixins: [modalMixin, modalInnerMixin],
  data() {
    return {
      tabs: [
        { component: 'CreateForm', name: this.$t('modals.createEntity.fields.form') },
        { component: 'ManageInfos', name: this.$t('modals.createEntity.fields.manageInfos') },
      ],
      currentComponent: 'CreateForm',
      showValidationErrors: true,
      form: {
        name: '',
        description: '',
        type: '',
        enabled: true,
        infos: [],
      },
    };
  },
  methods: {
    updateImpact(entities) {
      this.form.impacts = entities.map(entity => entity._id);
    },
    updateDependencies(entities) {
      this.form.dependencies = entities.map(entity => entity._id);
    },
    async create() {
      // TO DO
      // Entity creation
      // todo create object infos after entity creating
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();
      if (formIsValid) {
        await this.create();
      }
    },

  },
};
</script>

<style scoped>
  .tooltip {
    flex: 1 1 auto;
  }
</style>
