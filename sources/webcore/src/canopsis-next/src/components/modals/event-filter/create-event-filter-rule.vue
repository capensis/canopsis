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

import { eventFilterRuleToForm, formEnrichmentOptionsToEventFilterRule } from '@/helpers/forms/event-filter-rule';

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
      if (this.isEnrichmentType) {
        const isFormValid = await this.$validator.validateAll(['actions']);

        if (isFormValid) {
          await this.config.action({
            ...this.form.general,
            ...formEnrichmentOptionsToEventFilterRule(this.form.enrichmentOptions),
          });

          this.$modals.hide();
        }
      } else {
        await this.config.action({ ...this.form.general });

        this.$modals.hide();
      }
    },
  },
};
</script>

