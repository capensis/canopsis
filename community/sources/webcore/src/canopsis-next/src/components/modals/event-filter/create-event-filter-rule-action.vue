<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        event-filter-rule-action-form(v-model="form")
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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

import { eventFilterRuleActionToForm, formToEventFilterRuleAction } from '@/helpers/forms/event-filter-rule';

import EventFilterRuleActionForm from '@/components/other/event-filter/form/event-filter-rule-action-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilterRuleAction,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EventFilterRuleActionForm, ModalWrapper },
  mixins: [
    modalInnerMixin,
    submittableMixinCreator(),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: eventFilterRuleActionToForm(this.modal.config.ruleAction),
    };
  },
  computed: {
    title() {
      return this.modal.config.title || this.$t('modals.eventFilterRule.addAction');
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToEventFilterRuleAction(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
