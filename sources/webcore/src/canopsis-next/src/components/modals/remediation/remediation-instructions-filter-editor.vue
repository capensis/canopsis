<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ $t('common.filters') }}
      template(slot="text")
        v-layout(row)
          v-radio-group(
            v-model="form.with",
            hide-details
          )
            v-radio(
              label="With selected instructions",
              :value="true",
              color="primary"
            )
            v-radio(
              label="Without selected instructions",
              :value="false",
              color="primary"
            )
        v-layout(row)
          v-switch(
            v-model="form.all",
            label="Select all",
            color="primary",
            @change="changeSelectAll"
          )
        v-layout(row)
          v-select(
            v-model="form.instructions",
            v-validate="selectValidationRules",
            :items="remediationInstructions",
            :loading="remediationInstructionsPending",
            :disabled="form.all",
            label="Select instructions",
            :error-messages="errors.collect('instructions')",
            itemText="name",
            itemValue="name",
            name="instructions",
            multiple,
            clearable
          )
      template(slot="actions")
        v-btn(
          :disabled="submitting",
          depressed,
          flat,
          @click="$modals.hide"
        ) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesRemediationInstructionsMixin from '@/mixins/entities/remediation/instructions';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.remediationInstructionsFilterEditor,
  $_veeValidate: {
    validator: 'new',
  },
  components: { ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixin(),
    entitiesRemediationInstructionsMixin,
  ],
  data() {
    const defaultForm = {
      with: true,
      all: false,
      instructions: [],
    };

    const { filter } = this.modal.config;

    return {
      form: filter ? cloneDeep(filter) : defaultForm,
    };
  },
  computed: {
    selectValidationRules() {
      return this.form.all ? {} : { required: true };
    },
    actions() {
      return this.config.actions || {};
    },
  },
  mounted() {
    this.fetchRemediationInstructionsList({ limit: 10000 });
  },
  methods: {
    changeSelectAll(all) {
      if (all) {
        this.form.instructions = [];
      }
    },
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(this.form);
          }

          this.$modals.hide();
        } catch (err) {
          this.setFormErrors(err);
        }
      }
    },
  },
};
</script>
