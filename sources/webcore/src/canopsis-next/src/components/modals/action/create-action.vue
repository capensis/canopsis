<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        action-form(v-model="form", :disabledId="modal.config.item && !modal.config.isDuplicating")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(:disabled="isDisabled", type="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import { formToAction, actionToForm } from '@/helpers/forms/action';
import { generateActionId } from '@/helpers/entities';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import ActionForm from '@/components/other/action/form/action-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  inject: ['$system'],
  components: { ActionForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { item, isDuplicating } = this.modal.config;

    const form = actionToForm(item, this.$system.timezone);

    /**
     * If we're duplicating an action, generate a new unique id
     */
    if (isDuplicating) {
      form.generalParameters._id = generateActionId();
    }

    return {
      form,
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createAction.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          const data = formToAction(this.form, this.$system.timezone);

          await this.config.action(data);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
