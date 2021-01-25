<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper(close)
      template(slot="title")
        span {{ title }}
      template(slot="text")
        event-filter-form(
          v-model="form.general",
          :is-disabled-id-field="isDisabledIdField"
        )
        event-filter-enrichment-form(
          v-if="isEnrichmentType",
          v-model="form.enrichmentOptions"
        )
      template(slot="actions")
        v-btn(depressed, flat, @click="$modals.hide") {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { MODALS, EVENT_FILTER_RULE_TYPES } from '@/constants';

import { eventFilterRuleToForm, formToEventFilterRule } from '@/helpers/forms/event-filter-rule';

import authMixin from '@/mixins/auth';
import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';
import confirmableModalMixin from '@/mixins/confirmable-modal';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import EventFilterEnrichmentForm from '@/components/other/event-filter/form/event-filter-enrichment-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilterRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EventFilterForm, EventFilterEnrichmentForm, ModalWrapper },
  mixins: [
    authMixin,
    modalInnerMixin,
    submittableMixin(),
    confirmableModalMixin(),
  ],
  data() {
    const { rule } = this.modal.config;

    return {
      form: eventFilterRuleToForm(rule),
    };
  },
  computed: {
    title() {
      return this.config.title || this.$t('modals.eventFilterRule.create.title');
    },

    isDisabledIdField() {
      return this.config.isDisabledIdField;
    },

    isEnrichmentType() {
      return this.form.general.type === EVENT_FILTER_RULE_TYPES.enrichment;
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll(['actions']);

      if (isFormValid) {
        const eventFilter = formToEventFilterRule(this.form);
        eventFilter.author = this.currentUser._id;

        await this.config.action(eventFilter);

        this.$modals.hide();
      }
    },
  },
};
</script>

