<template>
  <v-select
    :value="value"
    :label="$t('modals.createPbehaviorType.fields.type')"
    :items="defaultTypes"
    :disabled="disabled"
    :loading="pending"
    item-value="type"
    item-text="name"
    return-object
    @input="updateValue"
  />
</template>

<script>
import { find } from 'lodash';

import { PBEHAVIOR_TYPE_TYPES } from '@/constants';

import { formBaseMixin } from '@/mixins/form';
import { entitiesPbehaviorTypeMixin } from '@/mixins/entities/pbehavior/types';

export default {
  mixins: [
    formBaseMixin,
    entitiesPbehaviorTypeMixin,
  ],
  props: {
    value: {
      type: String,
      default: PBEHAVIOR_TYPE_TYPES.active,
    },
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  data() {
    return {
      pending: false,
      defaultTypes: [],
    };
  },
  mounted() {
    this.fetchDefaultTypes();
  },
  methods: {
    async fetchDefaultTypes() {
      this.pending = true;

      const { data } = await this.fetchPbehaviorTypesListWithoutStore({
        params: { default: true },
      });

      this.defaultTypes = data.map(type => ({
        ...type,
        name: this.$t(`modals.createPbehaviorType.canonicalTypes.${type.type}`),
      }));

      this.pending = false;

      const activeType = find(data, { type: this.value });

      if (activeType) {
        this.$emit('update:color', activeType.color);
      }
    },

    updateValue(value = {}) {
      this.$emit('input', value.type, value.color);
    },
  },
};
</script>
