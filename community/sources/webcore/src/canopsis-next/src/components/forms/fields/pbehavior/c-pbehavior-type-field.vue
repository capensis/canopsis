<template lang="pug">
  v-select(
    v-field="value",
    v-validate="rules",
    :label="label || $t('common.type')",
    :loading="pending",
    :items="items",
    :error-messages="errors.collect(name)",
    :name="name",
    :disabled="disabled",
    :multiple="multiple",
    :chips="chips",
    :deletable-chips="chips",
    :small-chips="chips",
    :item-disabled="isItemDisabled",
    :return-object="returnObject",
    item-text="name",
    item-value="_id"
  )
</template>

<script>
import { MAX_LIMIT } from '@/constants';

import { mapIds } from '@/helpers/entities';

import { entitiesPbehaviorTypeMixin } from '@/mixins/entities/pbehavior/types';

export default {
  inject: ['$validator'],
  mixins: [entitiesPbehaviorTypeMixin],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: [Object, String, Array],
      default: '',
    },
    name: {
      type: String,
      default: 'type',
    },
    label: {
      type: String,
      default: '',
    },
    disabled: {
      type: Boolean,
      default: false,
    },
    multiple: {
      type: Boolean,
      default: false,
    },
    chips: {
      type: Boolean,
      default: false,
    },
    returnObject: {
      type: Boolean,
      default: false,
    },
    required: {
      type: Boolean,
      default: false,
    },
    max: {
      type: Number,
      required: false,
    },
    types: {
      type: Array,
      required: false,
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

      const { data: reasons } = await this.fetchPbehaviorTypesListWithoutStore({
        params: {
          types: this.types,
          limit: MAX_LIMIT,
        },
      });

      this.items = reasons;
      this.pending = false;
    },

    isItemDisabled(item) {
      if (this.max) {
        const types = this.returnObject
          ? mapIds(this.value)
          : this.value;

        return this.value.length === this.max && !types.includes(item._id);
      }

      return false;
    },
  },
};
</script>
