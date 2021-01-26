<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        pbehavior-form(v-model="form", :no-filter="noFilter")
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import { pbehaviorToForm, formToPbehavior, pbehaviorToRequest } from '@/helpers/forms/planning-pbehavior';

import authMixin from '@/mixins/auth';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import PbehaviorForm from '@/components/other/pbehavior/calendar/partials/pbehavior-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
  },
  components: { PbehaviorForm, ModalWrapper },
  mixins: [
    authMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { pbehavior, timezone } = this.modal.config;

    return {
      form: pbehaviorToForm(pbehavior, null, timezone),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createPbehavior.create.title');
    },

    noFilter() {
      return !!this.config.noFilter;
    },
  },
  methods: {
    async submit() {
      const isValid = await this.$validator.validateAll();

      if (isValid) {
        if (this.config.action) {
          await this.config.action(pbehaviorToRequest(formToPbehavior(this.form, this.config.timezone)));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
