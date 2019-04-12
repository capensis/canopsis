<template lang="pug">
  div
    snmp-rule-form-field-title(:label="$t('modals.createSnmpRule.fields.oid.title')")
    v-layout(row, wrap)
      v-flex.pr-1(xs6)
        v-autocomplete.pt-0(
        :value="form.moduleName",
        :items="modules",
        :search-input.sync="searchInput",
        :loading="modulesPending",
        :placeholder="$t('modals.createSnmpRule.fields.oid.labels.module')",
        hide-no-data,
        hide-details,
        @change="selectModuleAndFirstMib"
        )
      v-flex.pl-1(xs6)
        v-select.pt-0(
        :value="form.oid",
        :items="moduleMibs",
        :loading="moduleMibsPending",
        :menu-props="{ offsetY: true }",
        item-text="name",
        item-value="oid",
        return-object,
        hide-details,
        @input="selectModuleMib"
        )
</template>

<script>
import formMixin from '@/mixins/form';
import entitiesSnmpMibMixin from '@/mixins/entities/snmp-mib';

import SnmpRuleFormFieldTitle from './snmp-rule-form-field-title.vue';

export default {
  components: { SnmpRuleFormFieldTitle },
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
      await this.selectModule(this.form.moduleName, true);
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

    async selectModuleAndFirstMib(module) {
      const data = await this.selectModule(module);

      if (data && data.length) {
        this.selectModuleMib(data[0]);
      }
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

      return data;
    },

    selectModuleMib(moduleMib) {
      const model = { ...this.form };

      if (moduleMib) {
        model.mibName = moduleMib.name;
        model.oid = moduleMib.oid;
      } else {
        model.mibName = '';
        model.oid = '';
      }

      this.updateModel(model);
    },
  },
};
</script>
