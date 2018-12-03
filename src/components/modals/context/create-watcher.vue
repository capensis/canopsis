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
    v-layout.py-1(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click="submit") {{ $t('common.submit') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import modalInnerMixin from '@/mixins/modal/inner';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import { MODALS } from '@/constants';

import CreateForm from './partial/create-watcher-form.vue';
import ManageInfos from './partial/manage-infos.vue';

const { mapActions: watcherMapActions } = createNamespacedHelpers('watcher');

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
      tabs: [
        { component: 'CreateForm', name: this.$t('modals.createEntity.fields.form') },
        { component: 'ManageInfos', name: this.$t('modals.createEntity.fields.manageInfos') },
      ],
    };
  },
  methods: {
    ...watcherMapActions(['create', 'edit']),

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = {
          ...this.form,
          infos: this.form.infos,
          _id: this.config.item ? this.config.item._id : this.form.name,
          display_name: this.form.name,
          type: this.$constants.ENTITIES_TYPES.watcher,
          mfilter: this.form.mfilter,
        };

        try {
          if (this.config.item) {
            await this.edit({ data });
          } else {
            await this.create({ data });
          }

          this.refreshContextEntitiesLists();

          this.hideModal();
        } catch (err) {
          console.error(err);
        }
      }
    },
  },
};
</script>
