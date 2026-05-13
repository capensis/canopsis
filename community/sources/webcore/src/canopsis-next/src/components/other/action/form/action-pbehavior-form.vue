<template>
  <pbehavior-general-form
    v-field="value"
    :name-label="$t('common.namePrefix')"
    :name-tooltip="$t('scenario.tooltips.pbehaviorActionsNamePrefix')"
    :depth="1"
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
  },
};
</script>
