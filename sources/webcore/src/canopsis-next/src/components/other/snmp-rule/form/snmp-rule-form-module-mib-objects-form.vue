<template lang="pug">
  div
    snmp-rule-form-field(:label="label")
    v-layout(row, wrap)
      v-flex(xs12)
        v-menu(:items="items", full-width, offset-y, max-height="200")
          v-text-field.vars-input.pt-0(
            slot="activator",
            :value="form.value",
            :label="$t('modals.createSnmpRule.fields.moduleMibObjects.vars')",
            hide-details,
            @input="updateField('value', $event)"
          )
            template(slot="append", v-if="large")
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
            :value="form.regex",
            :label="$t('modals.createSnmpRule.fields.moduleMibObjects.regex')",
            hide-details,
            @input="updateField('regex', $event)"
          )
          v-text-field(
            :value="form.formatter",
            :label="$t('modals.createSnmpRule.fields.moduleMibObjects.formatter')",
            hide-details,
            @input="updateField('formatter', $event)"
          )
</template>

<script>
import formMixin from '@/mixins/form';

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

  .vars-input /deep/ .v-input__slot {
    height: 56px;
  }
</style>
