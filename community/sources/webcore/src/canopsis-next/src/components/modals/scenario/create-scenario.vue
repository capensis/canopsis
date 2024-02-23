<template>
  <v-form @submit.prevent="submit">
    <modal-wrapper close>
      <template #title="">
        <span>{{ title }}</span>
      </template>
      <template #text="">
        <scenario-form v-model="form" />
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
          :disabled="isDisabled"
          :loading="submitting"
          class="primary"
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

import { formToScenario, scenarioToForm, scenarioErrorToForm } from '@/helpers/entities/scenario/form';

import { modalInnerMixin } from '@/mixins/modal/inner';
import { entitiesScenarioMixin } from '@/mixins/entities/scenario';
import { validationErrorsMixinCreator } from '@/mixins/form/validation-errors';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import ScenarioForm from '@/components/other/scenario/form/scenario-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createScenario,
  $_veeValidate: {
    validator: 'new',
    delay: VALIDATION_DELAY,
  },
  inject: ['$system'],
  components: {
    ScenarioForm,
    ModalWrapper,
  },
  mixins: [
    modalInnerMixin,
    entitiesScenarioMixin,
    validationErrorsMixinCreator(),
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: scenarioToForm(this.modal.config.scenario, this.$system.timezone),
    };
  },
  computed: {
    title() {
      return this.config.title ?? this.$t('modals.createScenario.create.title');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            await this.config.action(formToScenario(this.form, this.$system.timezone));
          }

          this.$modals.hide();
        } catch (err) {
          if (err.error) {
            this.$popups.error({ text: err.error });
          } else {
            this.setFormErrors(scenarioErrorToForm(err, this.form));
          }
        }
      }
    },
  },
};
</script>
