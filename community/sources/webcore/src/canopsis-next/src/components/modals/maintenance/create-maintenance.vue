<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <c-alert type="warning">
          {{ config.warningText }}
        </c-alert>
        <maintenance-form v-model="form" />
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
          :loading="submitting"
          :disabled="isDisabled"
          class="primary"
          type="submit"
        >
          {{ $t(config.maintenance ? 'common.submit' : 'modals.createMaintenance.enableMaintenance') }}
        </v-btn>
      </template>
    </modal-wrapper>
  </v-form>
</template>

<script>
import { MODALS, VALIDATION_DELAY } from '@/constants';

import { maintenanceToForm } from '@/helpers/entities/maintenance/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import MaintenanceForm from '@/components/other/maintenance/form/maintenance-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createMaintenance,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { MaintenanceForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { maintenance } = this.modal.config;

    return {
      form: maintenanceToForm(maintenance),
    };
  },
  computed: {
    title() {
      return this.$t('modals.createMaintenance.setup.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validate();

      if (isFormValid) {
        await this.config.action?.(this.form);

        this.$modals.hide();
      }
    },
  },
};
</script>
