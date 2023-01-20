<template lang="pug">
  v-layout(row)
    template(v-if="value")
      v-text-field.mt-0.pt-0(:value="value", :disabled="disabled", readonly, hide-details)
      c-action-btn(
        :disabled="disabled",
        type="edit",
        btn-class="ml-2",
        @click="$emit('edit', value)"
      )
      c-action-btn(
        :disabled="disabled",
        type="delete",
        @click="$emit('remove')"
      )
    v-btn.ml-0(
      v-else,
      :disabled="disabled",
      :color="errors.has(name) ? 'error' : 'primary'",
      @click="$emit('add')"
    ) {{ $t('common.add') }}
</template>

<script>
export default {
  inject: ['$validator'],
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
    required: {
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

    disabled(value, oldValue) {
      if (value && !oldValue) {
        this.detachRequiredRule();
      } else {
        this.attachRequiredRule();
      }
    },
  },
  created() {
    if (!this.disabled) {
      this.attachRequiredRule();
    }
  },
  beforeDestroy() {
    this.detachRequiredRule();
  },
  methods: {
    attachRequiredRule() {
      const oldField = this.$validator.fields.find({ name: this.name });

      if (!oldField) {
        this.$validator.attach({
          name: this.name,
          rules: 'required:true',
          getter: () => this.value.length > 0,
          vm: this,
        });
      }
    },

    detachRequiredRule() {
      this.$validator.detach(this.name);
    },
  },
};
</script>
