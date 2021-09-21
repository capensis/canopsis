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

import { modalInnerMixin } from '@/mixins/modal/inner';
import { submittableMixinCreator } from '@/mixins/submittable';
import { confirmableModalMixinCreator } from '@/mixins/confirmable-modal';

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
    modalInnerMixin,
    submittableMixinCreator({
      formField: 'form.general',
    }),
    confirmableModalMixinCreator(),
  ],
  data() {
    return {
      form: eventFilterRuleToForm(this.modal.config.rule),
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
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(formToEventFilterRule(this.form));
        }

        this.$modals.hide();
      }
    },
  },
};
</script>
