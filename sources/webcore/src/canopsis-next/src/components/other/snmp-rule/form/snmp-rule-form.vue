<template lang="pug">
  div
    snmp-rule-form-field(label="oid")
      v-flex.pr-1(xs6)
        v-autocomplete.pt-0(
        placeholder="Select a mib module",
        :value="form.oid.moduleName",
        :items="modules",
        :search-input.sync="search",
        :loading="loading",
        hide-no-data,
        hide-details,
        @change="selectModule"
        )
      v-flex.pl-1(xs6)
        v-select.pt-0(
        :value="form.oid.mibName",
        :items="moduleMibs",
        item-text="name",
        loading,
        hide-details,
        offset-y,
        @input="updateField('oid.mibName', $event)"
        )
    v-layout(v-if="selectedModuleMib", row, wrap)
      v-alert(:value="selectedModuleMib", color="grey darken-1", outline) {{ selectedModuleMib.description }}
    snmp-rule-form-vars-field(
    :value="form.output",
    :items="selectedModuleMibObjects",
    label="output",
    large,
    @input="updateField('output', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.resource",
    :items="selectedModuleMibObjects",
    label="resource",
    large,
    @input="updateField('resource', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.component",
    :items="selectedModuleMibObjects",
    label="component",
    large,
    @input="updateField('component', $event)"
    )
    snmp-rule-form-vars-field(
    :value="form.connector_name",
    :items="selectedModuleMibObjects",
    label="connector_name",
    large,
    @input="updateField('connector_name', $event)"
    )
    snmp-rule-form-state-form(
    :form="form.state",
    :items="selectedModuleMibObjects",
    @input="updateField('state', $event)"
    )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import SnmpRuleFormStateForm from './snmp-rule-form-state-form.vue';
import SnmpRuleFormField from './snmp-rule-form-field.vue';
import SnmpRuleFormVarsField from './snmp-rule-form-vars-field.vue';

export default {
  components: { SnmpRuleFormStateForm, SnmpRuleFormField, SnmpRuleFormVarsField },
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

      modules: [],
      moduleMibs: [],
    };
  },
  computed: {
    selectedModuleMib() {
      return this.moduleMibs.find(({ name }) => name === this.form.oid.mibName);
    },

    selectedModuleMibObjects() {
      if (this.selectedModuleMib) {
        return Object.keys(this.selectedModuleMib.objects);
      }

      return [];
    },
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

      this.modules = data;

      this.loading = false;
    },
  },
  async mounted() {
    if (this.form.oid.moduleName) {
      await this.selectModule(this.form.oid.moduleName);
    }
  },
  methods: {
    async selectModule(module) {
      this.updateField('oid.moduleName', module);

      const { data } = await this.fetchSnmpMibList({
        params: {
          limit: 0,
          query: {
            nodetype: 'notification',
            moduleName: module,
          },
        },
      });

      this.moduleMibs = data;
    },
  },
};
</script>

<style lang="scss" scoped>
  .tt-uppercase {
    text-transform: uppercase;
  }
</style>
