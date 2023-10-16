<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.importPbehaviorException.title') }}</span>
      </template>
      <template #text="">
        <pbehavior-exception-import-form v-model="form" />
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

import { pbehaviorExceptionImportToForm } from '@/helpers/entities/pbehavior/exception/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';

import PbehaviorExceptionImportForm from '@/components/other/pbehavior/exceptions/form/pbehavior-exception-import-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.importPbehaviorException,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    PbehaviorExceptionImportForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
  ],
  data() {
    return {
      form: pbehaviorExceptionImportToForm(this.modal.config.pbehaviorException),
    };
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config?.action(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
