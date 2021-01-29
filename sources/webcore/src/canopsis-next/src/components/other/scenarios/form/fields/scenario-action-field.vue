<template lang="pug">
  v-card
    v-card-text
      v-layout(row)
        c-action-type-field(
          v-field="action.type",
          :disabled="disabled",
          :name="`${name}.type`",
          @input="errors.clear()"
        )
      v-layout(row)
        c-enabled-field(v-field="action.emit_trigger", :label="$t('scenario.fields.emitTrigger')")
      v-tabs(v-model="activeTab", centered, slider-color="primary")
        v-tab(:class="{ 'error--text': hasGeneralError }") {{ $t('scenario.tabs.general') }}
        v-tab(:class="{ 'error--text': hasPatternsError }") {{ $t('scenario.tabs.pattern') }}
      v-divider
      v-tabs-items.pt-2(v-model="activeTab")
        v-tab-item
          scenario-action-general-field(
            ref="general",
            v-field="action",
            :name="`${name}.parameters`"
          )
        v-tab-item
          scenario-action-pattern-field(ref="patterns", v-model="action.patterns", :name="name")
</template>

<script>
import formMixin from '@/mixins/form/object';

import ScenarioActionGeneralField from './scenario-action-general-field.vue';
import ScenarioActionPatternField from './scenario-action-pattern-field.vue';

export default {
  components: { ScenarioActionPatternField, ScenarioActionGeneralField },
  inject: ['$validator'],
  mixins: [formMixin],
  model: {
    prop: 'action',
    event: 'input',
  },
  props: {
    action: {
      type: Object,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'action',
    },
  },
  data() {
    return {
      activeTab: 0,
      hasGeneralError: false,
      hasPatternsError: false,
    };
  },
  mounted() {
    this.$watch(() => this.$refs.general.hasAnyError, (value) => {
      this.hasGeneralError = value;
    });

    this.$watch(() => this.$refs.patterns.hasAnyError, (value) => {
      this.hasPatternsError = value;
    });
  },
};
</script>
