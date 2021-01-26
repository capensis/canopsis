<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
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

import { generateEntityId } from '@/helpers/entities';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';

import EntityForm from '@/components/widgets/context/form/entity-form.vue';
import ManageInfos from '@/components/widgets/context/manage-infos.vue';

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
    entitiesContextEntityMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    return {
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
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const formData = { ...this.form };

        if (!this.config.item || this.config.isDuplicating) {
          formData._id = generateEntityId();
        }

        await this.config.action(formData);
        await this.refreshContextEntitiesLists();

        this.$modals.hide();
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
