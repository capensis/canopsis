<template>
  <v-layout>
    <template v-if="value">
      <v-text-field
        :value="value"
        :disabled="disabled"
        class="mt-0 pt-0"
        readonly
        hide-details
      />
      <c-action-btn
        :disabled="disabled"
        type="edit"
        btn-class="ml-2"
        @click="$emit('edit', value)"
      />
      <c-action-btn
        :disabled="disabled"
        type="delete"
        @click="$emit('remove')"
      />
    </template>
    <v-btn
      v-else
      :disabled="disabled"
      :color="errors.has(name) ? 'error' : 'primary'"
      class="ml-0"
      @click="$emit('add')"
    >
      {{ $t('common.add') }}
    </v-btn>
  </v-layout>
</template>

<script>
import { validationAttachRequiredMixin } from '@/mixins/form/validation-attach-required';

export default {
  inject: ['$validator'],
  mixins: [validationAttachRequiredMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      required: true,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    name: {
      type: String,
      default: 'storage',
    },
  },
  watch: {
    value() {
      this.$validator.validate(this.name);
    },

    disabled: {
      immediate: true,
      handler(disabled) {
        if (disabled) {
          this.detachRequiredRule();
        } else {
          this.attachRequiredRule(this.requiredRuleGetter);
        }
      },
    },
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    requiredRuleGetter() {
      return this.value.length > 0;
    },
  },
};
</script>
