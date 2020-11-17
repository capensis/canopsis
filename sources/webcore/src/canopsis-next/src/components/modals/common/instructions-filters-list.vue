<template lang="pug">
  modal-wrapper(data-test="filtersListModal")
    template(slot="title")
      span {{ $t('common.filters') }}
    template(slot="text")
      v-layout(row)
        v-radio-group(
          :value="condition",
          hide-details
        )
          v-radio(
            label="With selected instructions",
            :value="0",
            color="primary"
          )
          v-radio(
            label="Without selected instructions",
            :value="1",
            color="primary"
          )
      v-layout(row)
        v-switch(
          v-model="selectedAllInstructions",
          label="Select all",
          color="primary"
        )
      v-layout(row)
        v-select(
          v-model="selectedInstructions",
          :items="remediationInstructions",
          :loading="remediationInstructionsPending",
          :disabled="selectedAllInstructions",
          itemText="name",
          itemValue="_id",
          multiple
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
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import entitiesRemediationInstructionsMixin from '@/mixins/entities/remediation/instructions';

import ModalWrapper from '../modal-wrapper.vue';

/**
 * Confirmation modal
 */
export default {
  name: MODALS.filtersList,
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
    return {
      condition: 0,
      selectedInstructions: [],
      selectedAllInstructions: false,
    };
  },
  computed: {
    actions() {
      return this.config.actions || {};
    },
  },
  mounted() {
    this.fetchRemediationInstructionsList({ limit: 10000 });
  },
  methods: {
    toggle() {

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
