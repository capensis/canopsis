<template lang="pug">
  v-form(@submit.prevent="submit")
    modal-wrapper
      template(slot="title")
        span {{ config.title }}
      template(slot="text")
        event-filter-form(
          v-model="form",
          :isEditing="isEditing",
          :isDuplicating="isDuplicating"
        )
        event-filter-enrichment-form(
          v-if="form.type === $constants.EVENT_FILTER_RULE_TYPES.enrichment",
          v-model="enrichmentOptions"
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
import { cloneDeep } from 'lodash';
import { MODALS, EVENT_FILTER_RULE_TYPES, EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES } from '@/constants';

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
    return {
      ruleTypes: Object.values(EVENT_FILTER_RULE_TYPES),
      form: {
        _id: '',
        type: EVENT_FILTER_RULE_TYPES.drop,
        description: '',
        pattern: {},
        priority: 0,
        enabled: true,
      },
      enrichmentOptions: {
        actions: [],
        externalData: {},
        onSuccess: EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
        onFailure: EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
      },
    };
  },
  computed: {
    isEditing() {
      return !!this.config.rule;
    },
    isDuplicating() {
      return this.config.isDuplicating;
    },
  },
  mounted() {
    if (this.config.rule) {
      const {
        _id,
        type,
        description,
        pattern,
        priority,
        enabled = true,
        actions,
        external_data: externalData,
        on_success: onSuccess,
        on_failure: onFailure,
      } = cloneDeep(this.config.rule);

      this.form = {
        type,
        description,
        pattern,
        priority,
        enabled,
      };

      if (!this.isDuplicating) {
        this.form._id = _id;
      }

      this.enrichmentOptions = {
        actions,
        externalData,
        onSuccess: onSuccess || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
        onFailure: onFailure || EVENT_FILTER_ENRICHMENT_RULE_AFTER_TYPES.pass,
      };
    }
  },
  methods: {
    async submit() {
      if (this.form.type === EVENT_FILTER_RULE_TYPES.enrichment) {
        const isFormValid = await this.$validator.validateAll(['actions']);

        if (isFormValid) {
          await this.config.action({
            ...this.form,
            actions: this.enrichmentOptions.actions,
            external_data: this.enrichmentOptions.externalData,
            on_success: this.enrichmentOptions.onSuccess,
            on_failure: this.enrichmentOptions.onFailure,
          });

          this.$modals.hide();
        }
      } else {
        await this.config.action({ ...this.form });

        this.$modals.hide();
      }
    },
  },
};
</script>

