<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline {{ $t('modals.createAction.create.title') }}
    v-card-text
      action-form(v-model="form", :disableId="modal.config.item && !modal.config.isDuplicating")
    v-divider
    v-layout.py-1(justify-end)
      v-btn(depressed, flat, @click="hideModal") {{ $t('common.cancel') }}
      v-btn.primary(:disabled="errors.any()", @click="submit") {{ $t('common.actions.saveChanges') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import uuid from '@/helpers/uuid';
import { formToAction, actionToForm } from '@/helpers/forms/action';

import ActionForm from '@/components/other/action/form/action-form.vue';

export default {
  name: MODALS.createAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    ActionForm,
  },
  mixins: [modalInnerMixin],
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
          const data = formToAction({
            generalParameters: this.form.generalParameters,
            pbehaviorParameters: this.form.pbehaviorParameters,
            snoozeParameters: this.form.snoozeParameters,
          });

          await this.config.action(data);
        }

        this.hideModal();
      }
    },
  },
};
</script>
