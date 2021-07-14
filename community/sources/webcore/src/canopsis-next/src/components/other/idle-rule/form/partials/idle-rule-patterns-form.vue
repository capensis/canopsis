<template lang="pug">
  div
    patterns-list(
      v-if="isEntityType",
      v-field="form.entity_patterns",
      v-validate="'required'",
      name="entity_patterns",
      @input="change"
    )
    c-patterns-field(
      v-else,
      v-field="form",
      some-required,
      alarm,
      entity
    )
    v-alert(
      :value="true",
      type="warning",
      transition="fade-transition",
      dismissible
    )
      span Message
</template>

<script>
import { pick } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { formValidationHeaderMixin } from '@/mixins/form';

import PatternsList from '@/components/common/patterns-list/patterns-list.vue';

const { mapActions } = createNamespacedHelpers('idleRules');

export default {
  inject: ['$validator'],
  components: { PatternsList },
  mixins: [formValidationHeaderMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      default: () => ({}),
    },
    isEntityType: {
      type: Boolean,
      default: false,
    },
  },
  watch: {
    'form.entity_patterns': {
      handler() {
        console.log('CHANGED');
      },
    },
  },
  methods: {
    ...mapActions({
      fetchIdleRuleEntitiesCountWithoutStore: 'fetchEntitiesCountWithoutStore',
    }),

    async change() {
      const result = await this.fetchIdleRuleEntitiesCountWithoutStore({
        data: pick(this.form, ['entity_patterns', 'alarm_patterns']),
      });

      console.log(result);
    },
  },
};
</script>
