<template lang="pug">
  div
    v-layout(row, wrap)
      v-flex(xs12)
        .body-2 oid
      v-flex.pr-1(xs6)
        v-autocomplete.pt-0(
        placeholder="Select a mib module",
        :items="items",
        :search-input.sync="search",
        :loading="loading",
        hide-no-data,
        @input="selectModule"
        )
      v-flex.pl-1(xs6)
        v-select.pt-0(
        :items="anotherItems",
        item-text="name",
        return-object,
        offset-y,
        @input="selectSubModule"
        )
    snmp-rule-form-field(
    :value="form.output",
    :items="anotherAnotherItems",
    label="output",
    @input="updateField('output', $event)"
    )
    snmp-rule-form-field(
    :value="form.resource",
    label="resource",
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-field(
    :value="form.component",
    label="component",
    @input="updateField('component', $event)"
    )
    snmp-rule-form-field(
    :value="form.connector_name",
    label="connector_name",
    @input="updateField('connector_name', $event)"
    )
    v-layout(row, wrap)
      v-flex(xs12)
        .body-2 state
      v-flex(xs12)
        v-switch(label="To custom")
      v-flex(xs12)
        state-criticity-field(:value="form.state.state")
      v-flex(xs12)
        v-text-field(placeholder="Snmp vars match field")
      v-flex(xs12)
        v-layout(row, wrap)
          v-flex(xs2)
            v-chip(color="green", label)
              strong CRITICAL
          v-flex(xs10)
            v-text-field(placeholder="write template")
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormField from './snmp-rule-form-field.vue';

export default {
  components: { StateCriticityField, SnmpRuleFormField },
  mixins: [formMixin, entitiesSnmpMibMixin],
  model: {
    prop: 'form',
    event: 'input',
  },
  props: {
    form: {
      type: Object,
      required: true,
    },
  },
  data() {
    return {
      menuProps: {
        closeOnClick: false,
        closeOnContentClick: false,
        openOnClick: true,
      },
      search: '',
      items: [],
      anotherItems: [],
      anotherAnotherItems: [],
      loading: false,
    };
  },
  watch: {
    async search(value) {
      this.loading = true;

      const { data } = await this.fetchSnmpMibDistinctList({
        params: {
          limit: 10,
          projection: 'moduleName',
          query: {
            nodetype: 'notification',
            moduleName: { $regex: `.*${value || ''}.*`, $options: 'i' },
          },
        },
      });

      this.items = data;

      this.loading = false;
    },
  },
  methods: {
    async selectModule(module) {
      const { data } = await this.fetchSnmpMibList({
        params: {
          limit: 0,
          query: {
            nodetype: 'notification',
            moduleName: module,
          },
        },
      });

      this.anotherItems = data;
    },

    selectSubModule(subModule) {
      this.anotherAnotherItems = Object.keys(subModule.objects);
    },
  },
};
</script>
