<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :label="$t('modals.createPbehavior.steps.general.fields.type')",
    :loading="pbehaviorTypesPending",
    :items="pbehaviorTypes",
    :error-messages="errors.collect('type')",
    item-text="name",
    item-value="_id",
    name="type",
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
      required: false,
    },
  },
  mounted() {
    this.fetchPbehaviorTypesList();
  },
};
</script>
