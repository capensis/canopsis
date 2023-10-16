<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <pbehavior-form
          v-model="form"
          :no-pattern="noPattern"
        />
      </template>
      <template #actions="">
        <v-btn
          depressed
          text
          @click="$modals.hide"
        >
          {{ $t('common.cancel') }}
        </v-btn>
        <v-btn
          class="primary"
          :disabled="isDisabled"
          :loading="submitting"
          type="submit"
        >
          {{ $t('common.submit') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { pbehaviorToForm, formToPbehavior, pbehaviorToRequest } from '@/helpers/entities/pbehavior/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { authMixin } from '@/mixins/auth';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createPbehavior,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { PbehaviorForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    authMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
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

    noPattern() {
      return !!this.config.noPattern;
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
