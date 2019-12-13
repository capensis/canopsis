<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        watcher-form(v-model="form", :stack="stack")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { omit } from 'lodash';

import { MODALS, ENTITIES_TYPES, CANOPSIS_STACK } from '@/constants';

import uuid from '@/helpers/uuid';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesInfoMixin from '@/mixins/entities/info';

import WatcherForm from '@/components/other/context/form/watcher-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WatcherForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
    entitiesContextEntityMixin,
    entitiesInfoMixin,
  ],
  data() {
    const { item } = this.modal.config;

    let form = {
      name: '',
      mfilter: '{}',
      infos: {},
      impact: [],
      depends: [],
      entities: [],
      output_template: '',
    };

    if (item) {
      form = { ...item };
    }

    return {
      form,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        let data = {};

        if (this.stack === CANOPSIS_STACK.go) {
          data = {
            ...omit(this.form, ['mfilter', 'impact', 'depends']),
            _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uuid('watcher'),
            name: this.form.name,
            type: ENTITIES_TYPES.watcher,
            state: {
              method: 'worst',
            },
          };
        } else {
          data = {
            ...omit(this.form, ['entities', 'output_template']),
            _id: this.config.item && !this.config.isDuplicating ? this.config.item._id : uuid('watcher'),
            infos: this.form.infos,
            display_name: this.form.name,
            type: ENTITIES_TYPES.watcher,
          };
        }

        await this.config.action(data);
        await this.refreshContextEntitiesLists();

        this.$modals.hide();
      }
    },
  },
};
</script>

<style>

</style>
