<template lang="pug">
  v-menu(bottom, offset-y, :closeOnContentClick="false")
    v-btn.color-picker-btn(slot="activator", :style="buttonStyles", flat)
    chrome(value="value", @input="handleChange")
</template>

<script>
import { Chrome } from 'vue-color';
import { isUndefined, isObject } from 'lodash';

export default {
  components: { Chrome },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [String, Object],
      default: '#000',
    },
    returnType: {
      type: String,
      required: false,
    },
  },
  computed: {
    buttonStyles() {
      return { background: isObject(this.value) ? this.value.hex : this.value };
    },
  },
  methods: {
    handleChange(colors) {
      const color = colors[this.returnType];
      this.$emit('input', !isUndefined(color) ? color : colors);
    },
  },
};
</script>

<style lang="scss" scoped>
  .color-picker-btn {
    min-width: 30px;
    height: 30px;
    margin: 0;
    padding: 0;
    border: 1px solid black;
  }
</style>
