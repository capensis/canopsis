<template lang="pug">
  v-card
    v-card-title.primary.white--text
      v-layout(justify-space-between, align-center)
        span.headline Rule's pattern
    v-card-text
      v-tabs(fixed-tabs, v-model="activeTab")
        v-tab(v-for="(tab, key) in tabs", :key="key") {{ tab }}
      v-tabs-items(v-model="activeTab")
        v-tab-item
          pattern-field-creation-form
        v-tab-item
          pattern-simple-editor(v-model="config.pattern", :operators="operators")
        v-tab-item
          pattern-advanced-editor(v-model="config.pattern")
    v-divider
    v-layout.pa-2(justify-end)
      v-btn(@click="hideModal", depressed, flat) {{ $t('common.cancel') }}
      v-btn.primary(@click.prevent="submit") {{ $t('common.submit') }}
</template>

<script>
import { MODALS } from '@/constants';

import modalInnerMixin from '@/mixins/modal/modal-inner';

import PatternFieldCreationForm from './partial/pattern-field-creation-form.vue';
import PatternSimpleEditor from './partial/pattern-simple-editor.vue';
import PatternAdvancedEditor from './partial/pattern-advanced-editor.vue';

export default {
  name: MODALS.createEventFilterRulePattern,
  components: {
    PatternFieldCreationForm,
    PatternSimpleEditor,
    PatternAdvancedEditor,
  },
  mixins: [modalInnerMixin],
  data() {
    return {
      activeTab: 0,
      tabs: ['Add a field', 'Simple editor', 'Advanced editor'],
      pattern: {},
      operators: ['>=', '>', '<', '<=', 'regex'],
    };
  },
  methods: {
    submit() {
      this.config.action(this.config.pattern);
      this.hideModal();
    },
  },
};
</script>

