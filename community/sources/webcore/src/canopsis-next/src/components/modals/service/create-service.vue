<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ config.title }}</span>
      </template>
      <template #text="">
        <service-form v-model="form" />
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
          :disabled="isDisabled || advancedJsonWasChanged"
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
import { get } from 'lodash';

import { MODALS, VALIDATION_DELAY } from '@/constants';

import { serviceToForm, formToService } from '@/helpers/entities/service/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';
import { entitiesContextEntityMixin } from '@/mixins/entities/context-entity';

import ServiceForm from '@/components/other/service/form/service-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createService,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ServiceForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    entitiesContextEntityMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: serviceToForm(this.modal.config.item),
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
        await this.config.action(formToService(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
