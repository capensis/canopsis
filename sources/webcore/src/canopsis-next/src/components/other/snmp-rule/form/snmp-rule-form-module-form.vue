<template lang="pug">
  snmp-rule-form-field(label="oid")
    v-flex.pr-1(xs6)
      v-autocomplete.pt-0(
      placeholder="Select a mib module",
      :value="form.moduleName",
      :items="modules",
      :search-input.sync="searchInput",
      :loading="modulesPending",
      hide-no-data,
      hide-details,
      @change="selectModule"
      )
    v-flex.pl-1(xs6)
      v-select.pt-0(
      :value="form.mibName",
      :items="moduleMibs",
      :loading="moduleMibsPending",
      item-text="name",
      hide-details,
      offset-y,
      @input="updateField('mibName', $event)"
      )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import SnmpRuleFormField from './snmp-rule-form-field.vue';

export default {
  components: { SnmpRuleFormField },
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
    moduleMibs: {
      type: Array,
      default: () => [],
    },
  },
  data() {
    return {
      modules: [],
      searchInput: '',
      modulesPending: false,
      moduleMibsPending: false,
    };
  },
  watch: {
    searchInput(value) {
      this.fetchModulesList(value);
    },
  },
  async mounted() {
    if (this.form.moduleName) {
      await this.selectModule(this.form.moduleName);
    }
  },
  methods: {
    async fetchModulesList(searchInput) {
      this.modulesPending = true;

      const { data } = await this.fetchSnmpMibDistinctList({
        params: {
          limit: 10,
          projection: 'moduleName',
          query: {
            nodetype: 'notification',
            moduleName: { $regex: `.*${searchInput || ''}.*`, $options: 'i' },
          },
        },
      });

      this.modules = data;
      this.modulesPending = false;
    },

    async selectModule(module) {
      this.moduleMibsPending = true;

      this.updateField('moduleName', module);

      const { data } = await this.fetchSnmpMibList({
        params: {
          limit: 0,
          query: {
            nodetype: 'notification',
            moduleName: module,
          },
        },
      });

      this.$emit('update:moduleMibs', data);

      this.moduleMibsPending = false;
    },
  },
};
</script>
