<template lang="pug">
  v-select(
    v-field="value",
    v-validate="'required'",
    :label="$t('common.type')",
    :loading="pbehaviorTypesPending",
    :items="pbehaviorTypes",
    :error-messages="errorMessages",
    :name="name",
    :disabled="disabled",
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
    disabled: {
      type: Boolean,
      default: false,
    },
  },
  computed: {
    errorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.name, this.$t('common.type')));
    },
  },
  mounted() {
    this.fetchPbehaviorTypesList();
  },
};
</script>
