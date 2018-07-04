<template lang="pug">
  v-card
    v-card-title
      span.headline {{ $t('modals.createEntity.title') }}
    v-btn(
      v-for="tab in tabs",
      @click.prevent="currentComponent = tab.component",
      color="primary") {{ tab.name }}
    keep-alive
      component(
        :is="currentComponent",
        :name.sync="form.name",
        :description.sync="form.description",
        :type.sync="form.type",
        :enabled.sync="form.enabled",
        :infos.sync="form.infos",
      )
    v-text-field(
      :label="$t('common.name')",
      v-validate="'required|alpha'",
      v-model="lol",
      :error-messages="errorMsg",
      )
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
      lol: '',
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
  computed: {
    errorMsg() {
      return this.$validator.errors ? false : 'lol';
    },
  },
  methods: {
    async create() {
      // TO DO
      // Entity creation
    },
    async submit() {
      const formIsValid = await this.$validator.validateAll();
      console.log(formIsValid);
      console.log(this.$validator);
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
