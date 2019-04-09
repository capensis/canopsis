<template lang="pug">
  div
    snmp-rule-form-field(label="oid")
      v-flex.pr-1(xs6)
        v-autocomplete.pt-0(
        placeholder="Select a mib module",
        :items="items",
        :search-input.sync="search",
        :loading="loading",
        hide-no-data,
        hide-details,
        @input="selectModule"
        )
      v-flex.pl-1(xs6)
        v-select.pt-0(
        :items="anotherItems",
        item-text="name",
        return-object,
        hide-details,
        offset-y,
        @input="selectSubModule"
        )
    snmp-rule-form-vars-field(
    :value="form.output",
    :items="anotherAnotherItems",
    label="output",
    large,
    @input="updateField('output', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.resource",
    :items="anotherAnotherItems",
    label="resource",
    large,
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.component",
    :items="anotherAnotherItems",
    label="component",
    large,
    @input="updateField('component', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.connector_name",
    :items="anotherAnotherItems",
    label="connector_name",
    large,
    @input="updateField('connector_name', $event)"
    )
    snmp-rule-form-field(label="state")
      v-flex(xs12)
        v-switch(label="To custom")
    v-divider(light)
    v-layout(row, wrap)
      v-flex(xs12)
        state-criticity-field(:value="form.state.state")
    snmp-rule-form-vars-field(
    label="Define matching snmp var",
    :value="{}",
    :items="anotherAnotherItems",
    @input="updateField('state.stateoid', $event)"
    )
    v-layout(row, wrap)
      v-flex(xs12)
        v-layout(v-for="(state, key) in $constants.ENTITIES_STATES", :key="key", row, wrap, align-center)
          v-flex(xs2)
            v-chip(:style="{ backgroundColor: $constants.ENTITIES_STATES_STYLES[state].color }", label)
              strong.tt-uppercase {{ $t(`modals.createChangeStateEvent.states.${key}`) }}
          v-flex(xs10)
            v-text-field(
            :value="form.state[key]",
            placeholder="write template",
            @input="updateField(`state.${key}`, $event)"
            )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import StateCriticityField from '@/components/forms/fields/state-criticity-field.vue';

import SnmpRuleFormField from './snmp-rule-form-field.vue';
import SnmpRuleFormVarsField from './snmp-rule-form-vars-field.vue';

export default {
  components: { StateCriticityField, SnmpRuleFormField, SnmpRuleFormVarsField },
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

<style lang="scss" scoped>
  .tt-uppercase {
    text-transform: uppercase;
  }
</style>
