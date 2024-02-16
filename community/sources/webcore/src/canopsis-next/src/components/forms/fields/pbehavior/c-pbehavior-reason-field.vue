<template>
  <v-select
    v-field="reason"
    v-validate="rules"
    :label="$t('common.reason')"
    :loading="pending"
    :items="items"
    :error-messages="errors.collect(name)"
    :name="name"
    :return-object="returnObject"
    item-text="name"
    item-value="_id"
  />
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { entitiesPbehaviorReasonMixin } from '@/mixins/entities/pbehavior/reasons';

export default {
  inject: ['$validator'],
  mixins: [entitiesPbehaviorReasonMixin],
  model: {
    prop: 'reason',
    event: 'input',
  },
  props: {
    reason: {
      type: [Object, String],
      default: '',
    },
    name: {
      type: String,
      default: 'reason',
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
  mounted() {
    this.fetchList();
  },
  methods: {
    async fetchList() {
      this.pending = true;

      const { data: reasons } = await this.fetchPbehaviorReasonsListWithoutStore({
        params: { limit: MAX_LIMIT },
      });

      this.items = reasons;
      this.pending = false;
    },
  },
};
</script>
