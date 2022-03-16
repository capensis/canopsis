<template lang="pug">
  c-select-field(
    v-field="value",
    v-validate="rules",
    :label="$tc('common.pbehavior')",
    :loading="pending",
    :items="items",
    :error-messages="errors.collect(name)",
    :name="name",
    :return-object="returnObject",
    item-text="name",
    item-value="_id",
    autocomplete
  )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { entitiesPbehaviorMixin } from '@/mixins/entities/pbehavior';

export default {
  inject: ['$validator'],
  mixins: [entitiesPbehaviorMixin],
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
      default: 'pbehavior',
    },
    required: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      items: [],
      pending: false,
    };
  },
  computed: {
    rules() {
      return {
        required: this.required,
      };
    },
  },
  async mounted() {
    this.pending = true;

    const { data: pbehaviors } = await this.fetchPbehaviorsListWithoutStore({
      params: { limit: MAX_LIMIT },
    });

    this.items = pbehaviors;
    this.pending = false;
  },
};
</script>
