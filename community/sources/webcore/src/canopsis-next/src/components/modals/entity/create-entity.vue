<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <entity-form v-model="form" />
      </template>
      <template #actions="">
        <v-btn
          :disabled="submitting"
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

import { entityToForm, formToEntity } from '@/helpers/entities/entity/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import EntityForm from '@/components/other/entity/form/entity-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEntity,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: {
    EntityForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: entityToForm(this.modal.config.entity),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.createEntity.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        await this.config.action(formToEntity(this.form));

        this.$modals.hide();
      }
    },
  },
};
</script>
