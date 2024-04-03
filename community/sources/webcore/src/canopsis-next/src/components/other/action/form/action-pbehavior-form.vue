<template>
  <pbehavior-form
    v-field="value"
    no-pattern
    no-enabled
    no-comments
    with-start-on-trigger
  />
</template>

<script>
import PbehaviorForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-form.vue';

export default {
  components: { PbehaviorForm },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
  },
};
</script>
