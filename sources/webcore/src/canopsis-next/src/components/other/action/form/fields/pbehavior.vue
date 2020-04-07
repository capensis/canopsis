<template lang="pug">
  pbehavior-form(
    v-field="value",
    :author="$constants.ACTION_AUTHOR",
    noFilter
  )
</template>

<script>
import PbehaviorForm from '@/components/other/pbehavior/form/pbehavior-form.vue';

export default {
  inject: ['$validator'],
  components: {
    PbehaviorForm,
  },
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
