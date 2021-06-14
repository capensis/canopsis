<template lang="pug">
  v-autocomplete(
    v-field="value",
    :label="label",
    :hint="hint",
    :items="availableIconNames",
    :name="name",
    :error-messages="errors.collect(name)",
    persistent-hint
  )
    template(slot="selection", slot-scope="data")
      v-icon {{ data.item }}
      span.ml-2 {{ data.item }}
    template(slot="item", slot-scope="data")
      v-icon {{ data.item }}
      span.ml-2 {{ data.item }}
    template(slot="no-data")
      slot(name="no-data")
</template>

<script>
import materialIconNameByCode from '@/assets/material-icons/MaterialIcons-Regular.json';

export default {
  $_veeValidate: {
    name() {
      return this.name;
    },
    value() {
      return this.value;
    },
  },
  inject: ['$validator'],
  model: {
    prop: 'value',
    event: 'input',
  },
  props: {
    value: {
      type: String,
      default: '',
    },
    label: {
      type: String,
      default: '',
    },
    hint: {
      type: String,
      default: '',
    },
    name: {
      type: String,
      default: 'icon_name',
    },
  },
  computed: {
    availableIconNames() {
      return Object.keys(materialIconNameByCode);
    },
  },
};
</script>
