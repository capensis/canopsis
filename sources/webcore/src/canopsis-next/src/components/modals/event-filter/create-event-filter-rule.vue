<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ title }}
      template(slot="text")
        event-filter-form(
          v-model="form.general",
          :isDisabledIdField="isDisabledIdField"
        )
        event-filter-enrichment-form(
          v-if="isEnrichmentType",
          v-model="form.enrichmentOptions"
        )
      template(slot="actions")
        v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
        v-btn.primary(
          :disabled="isDisabled",
          :loading="submitting",
          type="submit"
        ) {{ $t('common.submit') }}
</template>

<script>
import { omit } from 'lodash';
import { MODALS, EVENT_FILTER_RULE_TYPES } from '@/constants';

import { eventFilterRuleToForm, formEnrichmentOptionsToEventFilterRule } from '@/helpers/forms/event-filter-rule';

import modalInnerMixin from '@/mixins/modal/inner';
import submittableMixin from '@/mixins/submittable';

import EventFilterForm from '@/components/other/event-filter/form/event-filter-form.vue';
import EventFilterEnrichmentForm from '@/components/other/event-filter/form/event-filter-enrichment-form.vue';

import ModalWrapper from '../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilterRule,
  $_veeValidate: {
    validator: 'new',
  },
  components: { EventFilterForm, EventFilterEnrichmentForm, ModalWrapper },
  mixins: [modalInnerMixin, submittableMixin()],
  data() {
    const { rule, isDuplicating } = this.modal.config;

    return {
      form: eventFilterRuleToForm(isDuplicating ? omit(rule, ['_id']) : rule),
    };
  },
  computed: {
    title() {
      let type = 'create';

      if (this.config.rule) {
        type = this.config.isDuplicating ? 'duplicate' : 'edit';
      }

      return this.$t(`modals.eventFilterRule.${type}.title`);
    },

    isDisabledIdField() {
      return this.config.rule && !this.config.isDuplicating;
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

