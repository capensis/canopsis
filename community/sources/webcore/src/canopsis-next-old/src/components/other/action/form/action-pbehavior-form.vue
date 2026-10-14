<template>
  <pbehavior-general-form
    v-field="value"
    :name="name"
    :name-label="$t('common.namePrefix')"
    :name-tooltip="$t('scenario.tooltips.pbehaviorActionsNamePrefix')"
    :depth="depth"
    block-child
    no-comments
    with-start-on-trigger
  />
</template>

<script>
import PbehaviorGeneralForm from '@/components/other/pbehavior/pbehaviors/form/pbehavior-general-form.vue';

export default {
  components: { PbehaviorGeneralForm },
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: Object,
      required: true,
    },
    name: {
      type: String,
      default: 'parameters',
    },
    depth: {
      type: Number,
      default: 0,
    },
  },
};
</script>
