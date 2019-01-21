<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ config.title }}
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
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit", :loading="submitting", :disabled="submitting") {{ $t('common.submit') }}
</template>

<script>
import { MODALS, ENTITIES_TYPES } from '@/constants';

import uuid from '@/helpers/uuid';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

import CreateForm from './partial/create-watcher-form.vue';
import ManageInfos from './partial/manage-infos.vue';

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    CreateForm,
    ManageInfos,
  },
  mixins: [modalInnerMixin, entitiesContextEntityMixin],
  data() {
    const { item } = this.modal.config;

    let form = {
      name: '',
      mfilter: '{}',
      infos: {},
      impact: [],
      depends: [],
    };

    if (item) {
      form = { ...item };
    }

    return {
      form,
      submitting: false,
      tabs: [
        { component: 'CreateForm', name: this.$t('modals.createEntity.fields.form') },
        { component: 'ManageInfos', name: this.$t('modals.createEntity.fields.manageInfos') },
      ],
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        this.submitting = true;

        const data = {
          ...this.form,
          _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uuid('watcher'),
          infos: this.form.infos,
          display_name: this.form.name,
          type: ENTITIES_TYPES.watcher,
        };

        try {
          await this.config.action(data);
          this.refreshContextEntitiesLists();

          this.hideModal();
        } catch (err) {
          this.addErrorPopup({ text: this.$t('error.default') });
          console.error(err);
        }

        this.submitting = false;
      }
    },
  },
};
</script>
