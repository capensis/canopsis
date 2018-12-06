<template lang="pug">
  div
</template>

<script>
export default {
  props: {
    value: {
      type: Object,
      required: true,
    },
    operators: {
      type: Array,
      required: true,
    },
  },
  methods: {
    isSimpleRule(rule) {
      if (typeof rule.value === 'string') {
        return true;
      }
      return false;
    },
  },
};
</script>
