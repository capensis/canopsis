<template lang="pug">
  v-select(
    v-field="value",
    v-validate="rules",
    :label="$t('common.type')",
    :loading="pbehaviorTypesPending",
    :items="types",
    :error-messages="errorMessages",
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
      type: [Object, String, Array],
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
    withIcon: {
      type: Boolean,
      default: false,
    },
    isItemDisabled: {
      type: Function,
      required: false,
    },
  },
  computed: {
    types() {
      return this.withIcon ? this.pbehaviorTypes.filter(type => type.icon_name) : this.pbehaviorTypes;
    },

    errorMessages() {
      return this.errors.collect(this.name).map(error => error.replace(this.name, this.$t('common.type')));
    },

    rules() {
      return {
        required: !this.disabled,
      };
    },
  },
  mounted() {
    this.fetchPbehaviorTypesList();
  },
};
</script>
