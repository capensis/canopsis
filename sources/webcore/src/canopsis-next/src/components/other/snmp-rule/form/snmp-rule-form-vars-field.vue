<template lang="pug">
  snmp-rule-form-field(:label="label")
    v-flex(xs12)
      v-menu(:items="items", full-width, offset-y, max-height="200")
        v-text-field.pt-0(
        slot="activator",
        :value="value.value",
        label="Snmp vars match field",
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
        :value="value.regex",
        label="Regex",
        hide-details,
        @input="updateField('regex', $event)"
        )
        v-text-field(
        :value="value.formatter",
        label="Format (capture group with \\x)",
        hide-details,
        @input="updateField('formatter', $event)"
        )
</template>

<script>
import formMixin from '@/mixins/form';

import SnmpRuleFormField from './snmp-rule-form-field.vue';

export default {
  components: { SnmpRuleFormField },
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
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
      isVisible: this.value.regex || this.value.format,
    };
  },
  methods: {
    toggleVisibility() {
      this.isVisible = !this.isVisible;
    },

    updateSelectableInput(item) {
      this.updateField('value', `${this.value.value || ''}{{ ${item} }}`);
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
</style>
