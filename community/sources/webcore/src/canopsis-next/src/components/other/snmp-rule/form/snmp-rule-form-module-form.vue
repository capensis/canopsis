<template>
  <div>
    <snmp-rule-form-field-title :label="$t('snmpRule.oid')" />
    <v-layout wrap>
      <v-flex
        class="pr-1"
        xs6
      >
        <v-autocomplete
          class="pt-0"
          v-validate="'required'"
          :value="form.moduleName"
          :items="modules"
          :search-input.sync="searchInput"
          :loading="modulesPending"
          :placeholder="$t('snmpRule.module')"
          :error-messages="errors.collect('moduleName')"
          item-text="moduleName"
          item-value="_id"
          name="moduleName"
          hide-no-data
          hide-details
          @change="selectModule"
        />
      </v-flex>
      <v-flex
        class="pl-1"
        xs6
      >
        <v-autocomplete
          class="pt-0"
          v-validate="'required'"
          :value="form.mib"
          :items="moduleMibs"
          :loading="moduleMibsPending"
          :menu-props="{ offsetY: true }"
          :error-messages="errors.collect('mib')"
          item-text="name"
          item-value="_id"
          name="mib"
          hide-no-data
          hide-details
          return-object
          @change="selectMib"
        />
      </v-flex>
    </v-layout>
  </div>
</template>

<script>
import { find, sortBy } from 'lodash';
import { createNamespacedHelpers } from 'vuex';

import { MAX_LIMIT } from '@/constants';

import { formMixin } from '@/mixins/form';

import SnmpRuleFormFieldTitle from './snmp-rule-form-field-title.vue';

const { mapActions } = createNamespacedHelpers('snmpMib');

export default {
  inject: ['$validator'],
  components: { SnmpRuleFormFieldTitle },
  mixins: [formMixin],
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
      modules: [],
      searchInput: '',
      modulesPending: false,
      moduleMibs: [],
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
    ...mapActions({
      fetchSnmpMibList: 'fetchList',
    }),

    async fetchModulesList(searchInput) {
      this.modulesPending = true;

      const { data } = await this.fetchSnmpMibList({
        params: {
          limit: MAX_LIMIT,
          nodetype: 'notification',
          search: searchInput,
          projection: 'moduleName',
          distinct: true,
        },
      });

      this.modules = sortBy(data, 'moduleName');
      this.modulesPending = false;
    },

    async selectModule(module) {
      this.moduleMibsPending = true;

      this.updateField('moduleName', module);

      const { data } = await this.fetchSnmpMibList({
        params: {
          limit: MAX_LIMIT,
          nodetype: 'notification',
          moduleName: module,
        },
      });

      this.moduleMibs = sortBy(data, 'name');
      this.moduleMibsPending = false;

      if (this.form.mib?.name) {
        this.updateField('mib', find(data, { name: this.form.mib.name }));
      }
    },

    selectMib(mib) {
      this.updateField('mib', mib);
    },
  },
};
</script>
