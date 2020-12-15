<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.eventFilterRule.editPattern') }}
    template(slot="text")
      pattern-form(
        v-model="form",
        :operators="operators",
        :is-simple-pattern="config.isSimplePattern"
      )
    template(slot="actions")
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(
        :disabled="patternWasChanged",
        @click.prevent="submit"
      ) {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep, get } from 'lodash';

import { MODALS, EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import PatternForm from '@/components/other/pattern/pattern-form.vue';

import ModalWrapper from '../../modal-wrapper.vue';

export default {
  name: MODALS.createEventFilterRulePattern,
  $_veeValidate: {
    validator: 'new',
  },
  components: {
    PatternForm,
    ModalWrapper,
  },
  mixins: [modalInnerMixin],
  data() {
    const { pattern = {} } = this.modal.config;

    return {
      form: cloneDeep(pattern),
      activeTab: 0,
    };
  },
  computed: {
    operators() {
      return this.config.operators || EVENT_FILTER_RULE_OPERATORS;
    },

    patternWasChanged() {
      return get(this.fields, ['pattern', 'touched']);
    },
  },
  methods: {
    async submit() {
      const isFormValid = await this.$validator.validateAll();

      if (isFormValid) {
        if (this.config.action) {
          await this.config.action(this.pattern);
        }

        this.$modals.hide();
      }
    },
  },
};
</script>

