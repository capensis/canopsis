<template lang="pug">
  v-layout.mb-4(row, wrap)
    v-flex(v-show="label", xs12)
      .body-2 {{ label }}
    v-flex(xs12)
      v-text-field.pt-0(
      :value="value.value",
      placeholder="Snmp vars match field",
      hide-details,
      @input="updateField('value', $event)"
      )
        template(slot="append")
          v-btn(
          :class="{ active: isVisible }",
          icon,
          @click="toggleVisibility"
          )
            v-icon attach_file
    v-expand-transition
      v-flex(v-show="isVisible", xs12)
        v-text-field(
        :value="value.regex",
        placeholder="Regex",
        hide-details,
        @input="updateField('regex', $event)"
        )
        v-text-field(
        :value="value.formatter",
        placeholder="Format (capture group with \\x)",
        hide-details,
        @input="updateField('formatter', $event)"
        )
</template>

<script>
import formMixin from '@/mixins/form';

export default {
  mixins: [formMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    label: {
      type: String,
      default: '',
    },
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      isVisible: false,
    };
  },
  methods: {
    toggleVisibility() {
      this.isVisible = !this.isVisible;
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
