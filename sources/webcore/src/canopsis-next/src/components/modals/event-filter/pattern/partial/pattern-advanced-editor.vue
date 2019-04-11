<template lang="pug">
  div
    v-textarea(:value="patternValue", @input="checkValidity", rows="15")
    v-btn(
    :color="isInvalidJSON ? 'error' : 'primary'",
    :disabled="isInvalidJSON",
    @click="save",
    ) {{ isInvalidJSON ? $t('errors.JSONNotValid') : $t('common.save') }}
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      newVal: {},
      isInvalidJSON: false,
      valueChanged: false,
    };
  },
  computed: {
    patternValue() {
      try {
        return JSON.stringify(this.value, undefined, 4);
      } catch (err) {
        console.error(err);
        return '';
      }
    },
  },
  methods: {
    checkValidity(value) {
      this.valueChanged = true;

      try {
        this.newVal = JSON.parse(value);
        this.isInvalidJSON = false;
      } catch (err) {
        this.isInvalidJSON = true;
      }
    },
    save() {
      if (this.valueChanged) {
        this.$emit('input', this.newVal);
      } else {
        this.$emit('input', this.value);
      }
    },
  },
};
</script>

