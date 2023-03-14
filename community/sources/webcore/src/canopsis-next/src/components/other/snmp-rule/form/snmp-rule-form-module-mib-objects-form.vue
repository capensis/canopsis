<template lang="pug">
  div
    snmp-rule-form-field(:label="label")
    v-layout(row, wrap)
      v-flex(xs12)
        v-menu(
          :items="items",
          max-height="200",
          full-width,
          offset-y
        )
          template(#activator="{ on }")
            v-text-field.vars-input.pt-0(
              v-on="on",
              v-field="form.value",
              :label="$t('snmpRule.moduleMibObjects')",
              hide-details
            )
              template(v-if="large", #append="")
                v-btn(
                  :class="{ active: isVisible }",
                  icon,
                  @click.stop="toggleVisibility"
                )
                  v-icon attach_file
          v-list
            v-list-tile(
              v-for="(item, index) in items",
              :key="index",
              @click="updateSelectableInput(item)"
            )
              v-list-tile-title {{ item }}
      v-expand-transition(v-if="large")
        v-flex(v-show="isVisible", xs12)
          v-text-field(
            v-field="form.regex",
            :label="$t('snmpRule.regex')",
            hide-details
          )
          v-text-field(
            v-field="form.formatter",
            :label="$t('snmpRule.formatter')",
            hide-details
          )
</template>

<script>
import { formMixin } from '@/mixins/form';

import SnmpRuleFormField from './snmp-rule-form-field-title.vue';

export default {
  components: { SnmpRuleFormField },
  mixins: [formMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
    items: {
      type: Array,
      default: () => [],
    },
    label: {
      type: String,
      default: '',
    },
    large: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      isVisible: this.form.regex || this.form.format,
    };
  },
  methods: {
    toggleVisibility() {
      this.isVisible = !this.isVisible;
    },

    updateSelectableInput(item) {
      this.updateField('value', `${this.form.value || ''}{{ ${item} }}`);
    },
  },
};
</script>

<style lang="scss" scoped>
  .v-btn.active {
    &:hover:before {
      opacity: .16;
    }

    &:before {
      background-color: currentColor;
    }
  }

  .vars-input ::v-deep .v-input__slot {
    height: 56px;
  }
</style>
