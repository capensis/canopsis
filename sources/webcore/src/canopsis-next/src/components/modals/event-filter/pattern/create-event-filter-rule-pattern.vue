<template lang="pug">
  modal-wrapper(close)
    template(slot="title")
      span {{ $t('modals.eventFilterRule.editPattern') }}
    template(slot="text")
      v-tabs(fixed-tabs, v-model="activeTab", slider-color="primary")
        v-tab(v-for="(tab, key) in tabs", :key="key") {{ tab }}
        v-tabs-items(v-model="activeTab")
          v-tab-item
            pattern-simple-editor(
              v-model="pattern",
              :operators="operators",
              :isSimplePattern="config.isSimplePattern"
            )
          v-tab-item
            pattern-advanced-editor(v-model="pattern")
    template(slot="actions")
      v-btn(@click="$modals.hide", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { cloneDeep } from 'lodash';

import { MODALS, EVENT_FILTER_RULE_OPERATORS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/inner';

import ModalWrapper from '../../modal-wrapper.vue';

import PatternSimpleEditor from './partial/pattern-simple-editor.vue';
import PatternAdvancedEditor from './partial/pattern-advanced-editor.vue';

export default {
  name: MODALS.createEventFilterRulePattern,
  components: {
    ModalWrapper,
    PatternSimpleEditor,
    PatternAdvancedEditor,
  },
  mixins: [modalInnerMixin],
  data() {
    const { pattern = {}, operators = EVENT_FILTER_RULE_OPERATORS } = this.modal.config;

    return {
      operators,

      pattern: cloneDeep(pattern),
      activeTab: 0,
      tabs: [
        this.$t('modals.eventFilterRule.simpleEditor'),
        this.$t('modals.eventFilterRule.advancedEditor'),
      ],
    };
  },
  methods: {
    async submit() {
      if (this.config.action) {
        await this.config.action(this.pattern);
      }

      this.$modals.hide();
    },
  },
};
</script>

