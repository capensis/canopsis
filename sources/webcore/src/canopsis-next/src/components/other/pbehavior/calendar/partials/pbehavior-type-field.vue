<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :label="$t('modals.createPbehavior.steps.general.fields.type')",
    :loading="pbehaviorTypesPending",
    :items="pbehaviorTypes",
    :error-messages="errors.collect(name)",
    :name="name",
    item-text="name",
    item-value="_id",
    return-object
  )
</template>

<script>
import entitiesPbehaviorTypesMixin from '@/mixins/entities/pbehavior/types';

export default {
  inject: ['$validator'],
  mixins: [entitiesPbehaviorTypesMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
  },
  mounted() {
    this.fetchPbehaviorTypesList();
  },
};
</script>
