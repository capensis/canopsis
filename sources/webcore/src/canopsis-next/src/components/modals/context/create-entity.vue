<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        v-tabs(slider-color="primary")
          v-tab(
            v-for="tab in tabs",
            :key="tab.name"
          ) {{ tab.name }}
          v-tab-item
            entity-form(v-model="form")
          v-tab-item
            manage-infos(v-model="form.infos")
      template(slot="actions")
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import uuid from '@/helpers/uuid';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

import EntityForm from '@/components/other/context/entity-form.vue';
import ManageInfos from '@/components/other/context/manage-infos.vue';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create an entity (watcher, resource, component, connector)
 */
export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    EntityForm,
    ManageInfos,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
    entitiesContextEntityMixin,
  ],
  data() {
    return {
      submitting: false,
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
    tabs() {
      return [
        { component: 'CreateForm', name: this.$t('modals.createEntity.fields.form') },
        { component: 'ManageInfos', name: this.$t('modals.createEntity.fields.manageInfos') },
      ];
    },
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

          this.$modals.hide();
        } catch (err) {
          console.error(err);
          this.$popups.error({ text: this.$t('error.default') });
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
