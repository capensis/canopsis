<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ $t('modals.createDynamicInfoInformation.create.title') }}</span>
      </template>
      <template #text="">
        <div>
          <v-text-field
            v-model="form.name"
            v-validate="nameRules"
            :label="$t('common.name')"
            :error-messages="errors.collect('name')"
            name
          />
          <c-mixed-field
            v-model="form.value"
            :label="$t('common.value')"
            name="value"
            required
          />
        </div>
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Modal to create dynamic info's information
 */
export default {
  name: MODALS.createDynamicInfoInformation,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    const { info = {} } = this.modal.config;

    return {
      form: {
        name: info.name ?? '',
        value: info.value ?? '',
      },
    };
  },
  computed: {
    initialName() {
      return this.config.info && this.config.info.name;
    },

    existingNames() {
      return this.config.existingNames;
    },

    nameRules() {
      return {
        required: true,
        unique: {
          values: this.existingNames,
          initialValue: this.initialName,
        },
      };
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.form);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
