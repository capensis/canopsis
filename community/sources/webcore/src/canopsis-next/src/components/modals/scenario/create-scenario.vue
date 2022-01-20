<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        scenario-form(v-model="form")
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

import { formToScenario, scenarioToForm, scenarioErrorToForm } from '@/helpers/forms/scenario';

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
      return this.config.title || this.$t('modals.createScenario.create.title');
    },

    originalPriority() {
      return this.config?.scenario?.priority;
    },
  },
  methods: {
    showConfirmScenarioPriorityChange(priority) {
      return new Promise((resolve) => {
        this.$modals.show({
          name: MODALS.confirmation,
          dialogProps: { persistent: true },
          config: {
            text: this.$t('scenario.errors.priorityExist', { priority }),
            action: () => {
              this.form.priority = priority;

              resolve();
            },
            cancel: resolve,
          },
        });
      });
    },

    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        try {
          if (this.config.action) {
            if (this.form.priority !== this.originalPriority) {
              const { valid, recommended_priority: recommendedPriority } = await this.checkScenarioPriority({
                data: { priority: this.form.priority },
              });

              if (!valid) {
                await this.showConfirmScenarioPriorityChange(recommendedPriority);
              }
            }

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
