<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        watcher-form(v-model="form", :stack="stack")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled || advancedJsonWasChanged",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { get } from 'lodash';

import { MODALS } from '@/constants';

import { watcherToForm, formToWatcher } from '@/helpers/forms/watcher';
import { generateWatcherId } from '@/helpers/entities';

import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';
import entitiesContextEntityMixin from '@/mixins/entities/context-entity';
import entitiesInfoMixin from '@/mixins/entities/info';

import WatcherForm from '@/components/widgets/context/form/watcher-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createWatcher,
  $_veeValidate: {
    validator: 'new',
  },
  components: { WatcherForm, ModalWrapper },
  mixins: [
    entitiesContextEntityMixin,
    entitiesInfoMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { item = {} } = this.modal.config;

    return {
      form: watcherToForm(item),
    };
  },
  computed: {
    advancedJsonWasChanged() {
      return get(this.fields, ['advancedJson', 'changed']);
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        const data = formToWatcher(this.form, this.stack);

        data._id = this.config.item && !this.config.isDuplicating
          ? this.config.item._id
          : generateWatcherId();

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
