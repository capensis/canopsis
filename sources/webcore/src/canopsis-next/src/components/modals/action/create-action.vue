<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('modals.createAction.create.title') }}
      template(slot="text")
        action-form(v-model="form", :disabledId="modal.config.item && !modal.config.isDuplicating")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import uuid from '@/helpers/uuid';
import { formToAction, actionToForm } from '@/helpers/forms/action';

import ActionForm from '@/components/other/action/form/action-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ActionForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { item, isDuplicating } = this.modal.config;

    const form = actionToForm(item);

    // If we're duplicating an action, generate a new unique id
    if (isDuplicating) {
      form.generalParameters._id = uuid('action');
    }

    return {
      form,
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const data = formToAction(this.form);

          await this.config.action(data);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
