<template lang="pug">
  v-flex
    v-btn.ml-0.btn-filter(
      :color="hasErrors ? 'error' : 'primary'",
      @click="showCreateFilterModal"
    ) {{ hasPattern ? $t('pbehavior.buttons.editFilter') : $t('pbehavior.buttons.addFilter') }}
    v-tooltip(v-if="hasPattern", fixed, top)
      template(#activator="{ on }")
        v-btn(v-on="on", icon)
          v-icon(color="grey darken-1") info
      span.pre {{ entityPattern | json }}
    v-alert(
      :value="hasErrors",
      type="error",
      transition="fade-transition"
    ) {{ errors.first(patternsFieldName) }}
</template>

<script>
import { MODALS } from '@/constants';

import { formGroupsToPatternRules } from '@/helpers/forms/pattern';

import { formMixin } from '@/mixins/form';

export default {
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'patterns',
    event: 'input',
  },
  props: {
    patterns: {
      type: Object,
      required: true,
    },
    patternsFieldName: {
      type: String,
      default: 'patterns',
    },
  },
  computed: {
    hasErrors() {
      return this.errors.has(this.patternsFieldName);
    },

    entityPattern() {
      return formGroupsToPatternRules(this.patterns.entity_pattern.groups);
    },

    hasPattern() {
      return this.entityPattern.length;
    },
  },
  watch: {
    patterns() {
      this.$validator.validate(this.patternsFieldName);
    },
  },
  created() {
    this.attachFilterRule();
  },
  beforeDestroy() {
    this.detachFilterRule();
  },
  methods: {
    attachFilterRule() {
      this.$validator.attach({
        name: this.patternsFieldName,
        rules: 'required:true',
        getter: () => !!this.hasPattern,
        vm: this,
      });
    },

    detachFilterRule() {
      this.$validator.detach(this.patternsFieldName);
    },

    showCreateFilterModal() {
      this.$modals.show({
        name: MODALS.pbehaviorPatterns,
        dialogProps: {
          zIndex: 300,
        },
        config: {
          patterns: this.patterns,
          withEntity: true,
          action: this.updateModel,
        },
      });
    },
  },
};
</script>

<style lang="scss" scoped>
.btn-filter.error {
  animation: shake .6s cubic-bezier(.25,.8,.5,1);
}

@keyframes shake {
  10%, 90% {
    transform: translate3d(-1px, 0, 0);
  }

  20%, 80% {
    transform: translate3d(2px, 0, 0);
  }

  30%, 50%, 70% {
    transform: translate3d(-4px, 0, 0);
  }

  40%, 60% {
    transform: translate3d(4px, 0, 0);
  }
}
</style>
